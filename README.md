# aquila 🦅

Implementation of the idea proposed by Puhua Guan in ["Cellular Automaton
Public-Key
Cryptosystem"](https://www.complex-systems.com/abstracts/v01_i01_a04/) paper.


## Roadmap

- [ ] Implement public/private key serialization/parsing
- [ ] Implement plain text padding
- [ ] Hash consing for symbolic math to reduce memory usage

## How it works?

Let *S* be a field of order m (`|S| = m`). The plain text is split into blocks
of *n* elements. For a public-key cryptosystem an invertible function is needed
that maps *S^n => S^n* and satisfies the following conditions:

* The function is easy to compute (for encryption)
* An inverse is hard to find (for decryption by an attacker)
* With some key information, an inverse can be easily computed

Invertible cellular automata with a large neighbourhood are good candidates for
that purpose.

### Invertible CA

Each block of a plain text can be encrypted by applying a sequence of *r*
invertible rules:

Each rule (referred as *MLISE* -- Multi-fold Linear Invertible System of
Equations -- in the codebase) is a set of *n* functions, which are constructed
by partitioning the variables into *s* "folds":

```
(x1, x2, ... xn) => (x_{1,1}, ...  x_{1,k1}), ... (x_{s,1}, ... x_{s,ks})
```


For each rule, the following conditions must hold for every fold:

1. **Noise Only from Previous Folds:** Each function in a part can include
   arbitrary expressions in variables from previous folds ("noise"), but not
   from the current or later ones. The first fold must not include noise.

2. **Partition isolation:** Variables from later parts must not appear in an
   earlier part's function.

3. **Linearity within a fold:** The variables within each fold must appear
   linearly in functions assigned to them.

4. **Invertible Coefficient Matrix:** For each fold, the linear coefficients
   (over *S*) form an invertible matrix.

#### Worked Example

Suppose `S` is GF(2) and block size `n = 5`.

Let's first partition the variables for each rule

```
# Rule 1 folds
(x1, x2, x3, x4, x5) => (x1, x2, x3), (x4, x5)

# Rule 2 folds
(x1', x2', x3', x4', x5') => (x4', x5'), (x1', x2', x3')
```

Now, let's put the linear part of the rules so that the variables of each fold
form an invertible matrix

```
# Rule 1:

x1' =    x2      \
x2' =       x3    | Fold 1
x3' = x1         /

x4' =    x5      \  Fold 2
x5' = x4         /

# Rule 2:

y1 = x4'         \  Fold 1
y2 =     x5'     /

y3 = x1'         \
y4 =     x2'      | Fold 2
y5 =         x3' /
```

Now add some "noise" in a variables from the previous folds

```
# Rule 1:

x1' =    x2
x2' =       x3
x3' = x1
x4' =    x5    + x1 * x2
x5' = x4       + x2 * x3
                 \_____/
                 "noise"
# Rule 2:

y1 = x4'
y2 =     x5'
y3 = x1'         + x4' * x5'
y4 =     x2'     + x1' * x4'
y5 =         x3' + x1'
                   \______/
                    "noise"
```

That's it! With that setup we can easily "encrypt" a block of plain text by
just evaluating the equations.

```
(y1, y2, ... yn) =  E_r(E_{r-1}(... E_1(x1, x2, ... xn)))
```

For finding the inverse we basically need to apply the inverse of each rule in
reversed order:

```
(x1, x2, ... xn) =  D_1(D_2(... D_r(y1, y2, ... yn)))
```

For extracting `D_i()` we'd better look at the equations in a vector form:

```
y1 =    x2
y2 =       x3
y3 = x1

y4 =    x5    + x1*x2
y5 = x4       + x2*x3
```

```
y1         =    x2
y2         =       x3
y3         = x1

y4 - x1*x2 =    x5
y5 - x2*x3 = x4
```

```
y_f1 = A * x_f1
y_f2 - n_f2 = A * x_f2
```

Now, the finding a solution is as easy as solving a system of linear equations
and adding the noise when necessary.

### General one-way CA, or making finding an inverse more difficult

We can represent our invertible CA in a form of non-linear system of equations
by composing the rules. That way, we can represent the `E()` as the following:


```
E(x1, x2, ... xn) =  E_r(E_{r-1}(... E_1(x1, x2, ... xn)))
```

For our example, the resulting rule composition may look like this:

```
y1 = x1*x2 + x5
y2 = x2*x3 + x4
y3 = x1*x2*x3 + x1*x2*x4 + x2*x3*x5 + x4*x5 + x2
y4 = x1*x2 + x2*x5 + x3
y5 = x1 + x2
```

Someone might say that's not hard to solve a system of equations like this on
paper, but as the neighborhood increase, it becomes infeasible to invert via
known algorithms, except with secret structure.
