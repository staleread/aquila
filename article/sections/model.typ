// Formal model of Aquila. Imported from the legacy note
// master/doc@b66acdc:master-article/notes/T001-model-and-derivation.typ.
// Labels named <h-…> are cited as evidence by hypotheses.yaml; do not rename them
// without updating the registry.

= Model and derivation <sec-model>

Let $m$, $s$, $f$, and $q$ be positive integers with $m=s dot f$. A cellular
automaton state is $X in FF_2^m$, and $q$ is the number of executed Rules. Each
Rule $R_t$, $0 <= t < q$, owns a permutation $pi_t$ of the coordinate indices
${0, dots, m-1}$ and folds $F_i$, $0 <= i < s$, of width $f$.

Each fold $F_i$ owns an invertible linear map $A_(t,i)$. The implementation
samples the off-diagonal entries of unit lower- and upper-triangular matrices
$L_(t,i)$ and $U_(t,i)$ and sets

$ A_(t,i) = L_(t,i) dot U_(t,i). $

Both factors, and hence $A_(t,i)$, are invertible over $FF_2$. Keeping the
factors permits inversion by forward and backward substitution in $O(f^2)$
time.

Folds ${F_i}_(i=1)^(s-1)$ also own a nonlinear component $Phi_(t,i)$ -- a
vector of $f$ boolean expressions in variables of predecessor local folds only.
For a Rule $t$ and fold $i$ let $x_0, dots, x_(i dot f-1)$ be the coordinates
of those predecessor folds. With a parameter $d$, where $2 <= d <= f$, to
control the nonlinear degree, output coordinate $ell$, $0 <= ell < f$, is

$ phi_ell^((d)) = xor.big_(r=2)^d product_(k=1)^r x_(alpha_(ell,r,k)), $

where $0 <= alpha_(ell,r,k) < i dot f$. For each monomial, its $r$ selected
indices are distinct. The complete nonlinear component is $Phi_(t,i)^((d)) =
(phi_0^((d)), dots, phi_(f-1)^((d)))$. Thus every output coordinate has one
monomial of each degree from $2$ through $d$. The parameter $d$ bounds the
degree of each generated confusion-map coordinate.

The permutation $pi_t$ maps a local coordinate $k$ to the corresponding global
state coordinate $pi_t (k)$. The resulting coordinate transport is:

$ G_(pi_t)(X)_k = X_(pi_t (k)), quad G_(pi_t)^(-1)(Y)_j = Y_(pi_t^(-1)(j)). $

Thus local fold $F_i, 0 <= i < s$ reads and writes precisely the global
coordinates $pi_t (f dot i+j)$ for $0 <= j < f$. The sets of coordinates of
distinct local folds are disjoint. Gathering a fold realizes $G_(pi_t)$ on that
fold; scattering it writes the local result back through $G_(pi_t)^(-1)$.

For each fold $F_i$, a vector $W_(t,i) in FF_2^f$ is gathered. The following
vector is then computed:

$ W'_(t,i) = A_(t,i)(W_(t,i)) xor Phi_(t,i)^((d))({W_(t,j)}_(j=0)^(i-1)), $

where $Phi_(t,0)^((d)) = arrow(0)$. The vector is then scattered. If $s=1$,
this is the only fold, so the Rule's local operation is simply the invertible
linear map $A_(t,0)$. For $s>1$, the remaining folds add the predecessor-only
nonlinear components defined above.

Together, these fold equations define the predecessor-triangular local map $T_t
(W_t)=W'_t$. Transporting it to global state coordinates defines the Rule

$ R_t = G_(pi_t)^(-1) compose T_t compose G_(pi_t). $

Operationally, applying a Rule $t$ preserves the snapshot of its input state
$X_t$ during the forward mapping. Since local fold coordinate sets are disjoint
and $Phi_(t,i)$ is evaluated against the input snapshot, the folds can be
processed in arbitrary order.

Consecutive Rules use different local coordinate systems. For $0 <= t < q-1$,
the relative transport from Rule $t$ to Rule $t+1$ is

$ H_(t+1,t) = G_(pi_(t+1)) compose G_(pi_t)^(-1), quad H_(t+1,t)(Y)_k =
Y_(pi_t^(-1)(pi_(t+1)(k))). $

It is the combined effect of scattering Rule $t$ and gathering Rule $t+1$;
accordingly,

$ R_(t+1) compose R_t = G_(pi_(t+1))^(-1) compose T_(t+1) compose H_(t+1,t)
compose T_t compose G_(pi_t). $

Forward mapping applies Rules in ascending index order and XORs the resulting
state with a non-zero affine shift vector $V$:

$ P(X_0) = V xor (R_(q-1) compose dots compose R_0)(X_0). $ <eq-public-map>

=== Invertibility <h-p-bijective>

For the inverse map, the $V$ shift is removed first. Rules are inverted in
descending index order. Each fold restores its local input vector $W_(t,i)$ by
XOR-ing away the noise and solving the resulting system of linear equations
using stored unit-LU factors and forward/backward substitutions:

$ W_(t,i) = A_(t,i)^(-1) (W'_(t,i) xor Phi_(t,i)^((d))({W_(t,j)}_(j=0)^(i-1))).
$

The result is then scattered back overwriting the values of the mapped global
coordinates with the restored input values.

Within $R_t^(-1)$, folds are restored sequentially in ascending order. The
first fold has no nonlinear term. For every later fold $i$, folds $0, dots,
i-1$ have already been restored, so $Phi_(t,i)$ noise can be recomputed and
removed. Applying $A_(t,i)^(-1)$ then uniquely restores $W_(t,i)$. Induction
over the folds proves that every Rule has a unique inverse. Removing $V$ and
applying the Rule inverses in descending order proves that $P$ is bijective.

=== Image of zero <h-p-zero>

Let $Q = R_(q-1) compose dots compose R_0$ be the map before the affine shift.
Every linear and nonlinear term in each Rule is constant-free, so $Q(0)=0$ and
therefore $P(0)=V$. This does not remove zero from the image: because $P$ is
bijective, exactly one input maps to zero.

