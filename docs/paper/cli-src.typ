#[
  = ДОДАТОК Б

  #v(-1.5em)
  #align(center)[*Програмний код консольного додатку*]
  #v(1em)
  
  #set par(
    first-line-indent: 0em,
    leading: 0.7em
  )

  #columns(2)[
  
  Файл `analyze.go`

  ```
  package main

  import (
    "crypto/rand"
    "encoding/csv"
    "fmt"
    "io"
    "math/bits"
    "os"
    "path/filepath"
    "strconv"

    "github.com/staleread/aquila/asym"
  )

  func runAnalyze(folder string, rndSamples int, correlation bool) error {
    priv, err := asym.GeneratePrivateKey(rand.Reader)
    if err != nil {
      return fmt.Errorf("failed to generate private key: %w", err)
    }

    blockSizeBytes := asym.BlockSize
    blockSizeBits := blockSizeBytes * 8

    changeCounts := make([]int, blockSizeBits)

    if err := os.MkdirAll(folder, 0755); err != nil {
      return fmt.Errorf("failed to create directory %s: %w", folder, err)
    }

    avFile, err := os.Create(filepath.Join(folder, "avalanche.csv"))
    if err != nil {
      return fmt.Errorf("failed to create avalanche.csv: %w", err)
    }
    defer avFile.Close()

    writer := csv.NewWriter(avFile)
    defer writer.Flush()

    if err := writer.Write([]string{"bit_index", "zeros_bg", "ones_bg", "random_bg_avg"}); err != nil {
      return fmt.Errorf("failed to write CSV header: %w", err)
    }

    hammingDistance := func(a, b []byte) int {
      dist := 0
      for i := range a {
        dist += bits.OnesCount8(a[i] ^ b[i])
      }
      return dist
    }

    zeroBaseSrc := make([]byte, blockSizeBytes)
    zeroBaseDst, err := priv.Sign(nil, zeroBaseSrc, nil)
    if err != nil {
      return fmt.Errorf("failed to encrypt zero base: %w", err)
    }

    oneBaseSrc := make([]byte, blockSizeBytes)
    for j := range oneBaseSrc {
      oneBaseSrc[j] = 0xFF
    }
    oneBaseDst, err := priv.Sign(nil, oneBaseSrc, nil)
    if err != nil {
      return fmt.Errorf("failed to encrypt one base: %w", err)
    }

    randSamplesSrc := make([][]byte, rndSamples)
    randSamplesDst := make([][]byte, rndSamples)
    for s := range rndSamples {
      randSamplesSrc[s] = make([]byte, blockSizeBytes)
      if _, err := io.ReadFull(rand.Reader, randSamplesSrc[s]); err != nil {
        return fmt.Errorf("failed to read random bytes: %w", err)
      }
      randSamplesDst[s], err = priv.Sign(nil, randSamplesSrc[s], nil)
      if err != nil {
        return fmt.Errorf("failed to encrypt random sample: %w", err)
      }
    }

    srcBuf := make([]byte, blockSizeBytes)
    var dstBuf []byte

    for i := range blockSizeBits {
      // --- Zeros BG ---
      for j := range srcBuf {
        srcBuf[j] = 0
      }
      srcBuf[i/8] = 1 << (i % 8)
      dstBuf, err = priv.Sign(nil, srcBuf, nil)
      if err != nil {
        return fmt.Errorf("failed to encrypt zeros bg: %w", err)
      }
      zerosBgVal := hammingDistance(dstBuf, zeroBaseDst)

      // --- Ones Background ---
      for j := range srcBuf {
        srcBuf[j] = 0xFF
      }
      srcBuf[i/8] &^= (1 << (i % 8))
      dstBuf, err = priv.Sign(nil, srcBuf, nil)
      if err != nil {
        return fmt.Errorf("failed to encrypt ones bg: %w", err)
      }
      onesBgVal := hammingDistance(dstBuf, oneBaseDst)

      // --- Random Background Avg ---
      totalRandDiff := 0
      for s := range rndSamples {
        copy(srcBuf, randSamplesSrc[s])
        srcBuf[i/8] ^= (1 << (i % 8))
        dstBuf, err = priv.Sign(nil, srcBuf, nil)
        if err != nil {
          return fmt.Errorf("failed to encrypt random background: %w", err)
        }
        totalRandDiff += hammingDistance(dstBuf, randSamplesDst[s])

        for o := range blockSizeBits {
          if ((randSamplesDst[s][o/8] ^ dstBuf[o/8]) & (1 << (o % 8))) != 0 {
            changeCounts[o]++
          }
        }
      }
      randBgAvgVal := float64(totalRandDiff) / float64(rndSamples)

      row := []string{
        strconv.Itoa(i),
        strconv.Itoa(zerosBgVal),
        strconv.Itoa(onesBgVal),
        fmt.Sprintf("%.4f", randBgAvgVal),
      }
      if err := writer.Write(row); err != nil {
        return fmt.Errorf("failed to write row %d to avalanche.csv: %w", i, err)
      }
    }

    if correlation {
      pub, err := priv.PublicKey()
      if err != nil {
        return fmt.Errorf("failed to derive public key: %w", err)
      }
      desc := pub.Describe()

      corrFile, err := os.Create(filepath.Join(folder, "monom-correlation.csv"))
      if err != nil {
        return fmt.Errorf("failed to create monom-correlation.csv: %w", err)
      }
      defer corrFile.Close()

      corrWriter := csv.NewWriter(corrFile)
      defer corrWriter.Flush()

      if err := corrWriter.Write([]string{"output_bit", "monomial_count", "avalanche_prob"}); err != nil {
        return fmt.Errorf("failed to write CSV header to monom-correlation.csv: %w", err)
      }

      for o := range blockSizeBits {
        prob := float64(changeCounts[o]) / float64(blockSizeBits*rndSamples)
        row := []string{
          strconv.Itoa(o),
          strconv.Itoa(desc.MonomialCounts[o]),
          fmt.Sprintf("%.6f", prob),
        }
        if err := corrWriter.Write(row); err != nil {
          return fmt.Errorf("failed to write row %d to monom-correlation.csv: %w", o, err)
        }
      }
    }

    return nil
  }
  ```

  Файл `anf.go`

  ```
  package main

  import (
    "bufio"
    "crypto/rand"
    "fmt"
    "os"

    "github.com/staleread/aquila/asym"
  )

  func exportANF(keyPath string) error {
    keyF, err := os.Open(keyPath)
    if err != nil {
      return fmt.Errorf("failed to open public key file: %w", err)
    }
    defer keyF.Close()

    pub := &asym.PublicKey{}
    if err := pub.Decode(keyF); err != nil {
      return fmt.Errorf("failed to decode public key: %w", err)
    }

    randInput := make([]byte, asym.BlockSize)
    if _, err := rand.Read(randInput); err != nil {
      return fmt.Errorf("failed to generate random input block: %w", err)
    }

    writer := bufio.NewWriter(os.Stdout)
    if err := pub.ExportToANF(writer, randInput); err != nil {
      return fmt.Errorf("failed to export ANF: %w", err)
    }

    if err := writer.Flush(); err != nil {
      return fmt.Errorf("failed to flush buffer: %w", err)
    }

    return nil
  }
  ```

  Файл `build.go`

  ```
  package main

  import (
    "fmt"
    "os"
    "os/exec"
  )

  func runBuild(configID string) error {
    var block, comp, fold, deg int
    n, err := fmt.Sscanf(configID, "%dc%df%dd%d", &block, &comp, &fold, &deg)
    if err != nil || n != 4 {
      return fmt.Errorf("invalid config ID format (must be <block>c<comp>f<fold>d<deg>): %w", err)
    }

    tags := fmt.Sprintf("block%d comp%d fold%d deg%d", block, comp, fold, deg)
    fmt.Printf("Building aquila-cli with tags: %s\n", tags)

    cmd := exec.Command("go", "build", "-tags", tags)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
      return fmt.Errorf("failed to run go build: %w", err)
    }

    fmt.Println("Build successful!")
    return nil
  }
  ```

  Файл `config.go`

  ```
  package main

  import (
    "fmt"

    "github.com/staleread/aquila/asym"
  )

  func showConfig() {
    fmt.Printf("Block size: %d\n", asym.BlockSize*8)
    fmt.Printf("Compositions: %d\n", asym.Compositions)
    fmt.Printf("Fold size: %d\n", asym.FoldSize)
    fmt.Printf("Confusion Degree: %d\n", asym.Degree)
  }

  func getConfigID() string {
    return fmt.Sprintf("%dc%df%dd%d", asym.BlockSize*8, asym.Compositions, asym.FoldSize, asym.Degree)
  }
  ```

  Файл `decrypt.go`

  ```
  package main

  import (
    "crypto/rand"
    "fmt"
    "os"

    "github.com/staleread/aquila/asym"
  )

  func decryptFile(inputPath, outputPath, keyPath string) error {
    keyF, err := os.Open(keyPath)
    if err != nil {
      return fmt.Errorf("failed to open private key file: %w", err)
    }
    defer keyF.Close()

    priv := &asym.PrivateKey{}
    if err := priv.Decode(keyF); err != nil {
      return fmt.Errorf("failed to decode private key: %w", err)
    }

    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
      return fmt.Errorf("failed to read input file: %w", err)
    }

    if len(ciphertext)%asym.BlockSize != 0 {
      return fmt.Errorf("ciphertext length is not a multiple of the block size")
    }

    paddedPlaintext, err := priv.Decrypt(rand.Reader, ciphertext, nil)
    if err != nil {
      return fmt.Errorf("failed to decrypt data: %w", err)
    }

    plaintext, err := pkcs7Unpad(paddedPlaintext, asym.BlockSize)
    if err != nil {
      return fmt.Errorf("failed to unpad decrypted data: %w", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
      return fmt.Errorf("failed to write output file: %w", err)
    }

    fmt.Printf("Successfully decrypted %s to %s\n", inputPath, outputPath)
    return nil
  }

  func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
    length := len(data)
    if length == 0 {
      return nil, fmt.Errorf("invalid padding size")
    }

    unpadding := int(data[length-1])
    if unpadding > blockSize || unpadding == 0 {
      return nil, fmt.Errorf("invalid padding")
    }

    padtext := data[length-unpadding:]
    for _, b := range padtext {
      if int(b) != unpadding {
        return nil, fmt.Errorf("invalid padding")
      }
    }
    return data[:(length - unpadding)], nil
  }
  ```

  Файл `encrypt.go`

  ```
  package main

  import (
    "bytes"
    "crypto/rand"
    "fmt"
    "os"

    "github.com/staleread/aquila/asym"
  )

  func encryptFile(inputPath, outputPath, keyPath string) error {
    keyF, err := os.Open(keyPath)
    if err != nil {
      return fmt.Errorf("failed to open public key file: %w", err)
    }
    defer keyF.Close()

    pub := &asym.PublicKey{}
    if err := pub.Decode(keyF); err != nil {
      return fmt.Errorf("failed to decode public key: %w", err)
    }

    inputData, err := os.ReadFile(inputPath)
    if err != nil {
      return fmt.Errorf("failed to read input file: %w", err)
    }

    paddedData := pkcs7Pad(inputData, asym.BlockSize)

    ciphertext, err := pub.Encrypt(rand.Reader, paddedData)
    if err != nil {
      return fmt.Errorf("failed to encrypt data: %w", err)
    }

    if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
      return fmt.Errorf("failed to write output file: %w", err)
    }

    fmt.Printf("Successfully encrypted %s to %s\n", inputPath, outputPath)
    return nil
  }

  func pkcs7Pad(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padtext...)
  }
  ```

  Файл `generate.go`

  ```
  package main

  import (
    "crypto/rand"
    "fmt"
    "os"

    "github.com/staleread/aquila/asym"
  )

  func generateKeyPair(name string) error {
    cfgID := getConfigID()
    privFile := fmt.Sprintf("id_aquila%s", cfgID)
    pubFile := fmt.Sprintf("id_aquila%s.pub", cfgID)

    if name != "" {
      privFile = fmt.Sprintf("id_aquila%s_%s", cfgID, name)
      pubFile = fmt.Sprintf("id_aquila%s_%s.pub", cfgID, name)
    }

    priv, pub, err := asym.GenerateKeyPair(rand.Reader)
    if err != nil {
      return fmt.Errorf("failed to generate key pair: %w", err)
    }

    privF, err := os.OpenFile(privFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
      return fmt.Errorf("failed to open private key file: %w", err)
    }
    defer privF.Close()

    if err := priv.Encode(privF); err != nil {
      return fmt.Errorf("failed to write private key: %w", err)
    }

    pubF, err := os.OpenFile(pubFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
      return fmt.Errorf("failed to open public key file: %w", err)
    }
    defer pubF.Close()

    if err := pub.Encode(pubF); err != nil {
      return fmt.Errorf("failed to write public key: %w", err)
    }

    fmt.Printf("Successfully generated key pair:\n  Private key: %s\n  Public key:  %s\n", privFile, pubFile)
    return nil
  }
  ```

  Файл `main.go`

  ```
  package main

  import (
    "fmt"
    "os"

    "github.com/alecthomas/kong"
  )

  var CLI struct {
    Gen struct {
      Name string `short:"n" help:"Optional name for the key pair."`
    } `cmd:"" help:"Generate a new key pair."`

    Enc struct {
      Input  string `short:"i" required:"" type:"existingfile" help:"Path to the input file."`
      Output string `short:"o" required:"" type:"path" help:"Path to the output file."`
      Key    string `short:"k" required:"" type:"existingfile" help:"Path to the public key file."`
    } `cmd:"" help:"Encrypt a file."`

    Dec struct {
      Input  string `short:"i" required:"" type:"existingfile" help:"Path to the input file."`
      Output string `short:"o" required:"" type:"path" help:"Path to the output file."`
      Key    string `short:"k" required:"" type:"existingfile" help:"Path to the private key file."`
    } `cmd:"" help:"Decrypt a file."`

    Config struct{} `cmd:"" help:"Show cipher configuration."`

    Build struct {
      ConfigID string `arg:"" help:"Configuration ID in format <block>c<comp>f<fold>d<deg>."`
    } `cmd:"" help:"Build the CLI with the specified cipher configuration."`

    Analyze struct {
      Folder      string `short:"o" required:"" type:"path" help:"Path to the output folder."`
      RndSamples  int    `long:"rnd-samples" required:"" help:"Number of random samples to analyze."`
      Correlation bool   `long:"correlation" help:"Generate monom-correlation.csv using public key (caution: deriving public key for large configurations like 96-bit key with 2 compositions can result in ~3 GB public key)."`
    } `cmd:"" name:"analyze" help:"Run cipher analysis experiments."`

    Anf struct {
      Key string `short:"k" required:"" type:"existingfile" help:"Path to the public key file."`
    } `cmd:"" help:"Export public key equations in Algebraic Normal Form (ANF)."`
  }

  func main() {
    ctx := kong.Parse(&CLI)
    switch ctx.Command() {
    case "gen":
      if err := generateKeyPair(CLI.Gen.Name); err != nil {
        fmt.Fprintf(os.Stderr, "Error generating keys: %v\n", err)
        os.Exit(1)
      }
    case "enc":
      if err := encryptFile(CLI.Enc.Input, CLI.Enc.Output, CLI.Enc.Key); err != nil {
        fmt.Fprintf(os.Stderr, "Error encrypting file: %v\n", err)
        os.Exit(1)
      }
    case "dec":
      if err := decryptFile(CLI.Dec.Input, CLI.Dec.Output, CLI.Dec.Key); err != nil {
        fmt.Fprintf(os.Stderr, "Error decrypting file: %v\n", err)
        os.Exit(1)
      }
    case "config":
      showConfig()
    case "build <config-id>":
      if err := runBuild(CLI.Build.ConfigID); err != nil {
        fmt.Fprintf(os.Stderr, "Error building CLI: %v\n", err)
        os.Exit(1)
      }
    case "analyze":
      if err := runAnalyze(CLI.Analyze.Folder, CLI.Analyze.RndSamples, CLI.Analyze.Correlation); err != nil {
        fmt.Fprintf(os.Stderr, "Error running analysis: %v\n", err)
        os.Exit(1)
      }
    case "anf":
      if err := exportANF(CLI.Anf.Key); err != nil {
        fmt.Fprintf(os.Stderr, "Error exporting ANF: %v\n", err)
        os.Exit(1)
      }
    }
  }
  ```
    
  ]
]