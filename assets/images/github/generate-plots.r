# plotting function
create_plot <- function(name, xlab, ylab, f, a) {
  xs <- lapply(c(0:100), function(x) x/100)
  ys <- lapply(xs, function(y) f(y,a))
  
  heading = name
  plot(xs, ys, type="n", main=heading, xlab=xlab, ylab=ylab)
  lines(xs, ys, type="l")
}
# plotting function
create_plot3D <- function(name, xlab, ylab, f, k) {
  xs <- lapply(c(0:100), function(x) x/100)
  ys <- lapply(c(0:100), function(x) x/100)
  zs <- outer(xs, ys, function(x, y) f(x,y,k))
  
  persp(xs, ys, zs,
        main="Perspective Plot of a Cone",
        zlab = "Height",
        theta = 30, phi = 15,
        col = "springgreen", shade = 0.5)
}

# normal distribution function
d <- function(nDoth, alpha) {
  (alpha^2) / (pi * (nDoth^2 * (alpha^2 - 1) + 1)^2)
}

# geometry function
geomSmith <- function(nDotv, k) {
  nDotv / (nDotv * (1-k) + k)
}
g <- function(nDotv, nDotl, k) {
  geomSmith(nDotv, k) * geomSmith(nDotl, k)
}

alpha <- 0.1
create_plot("Normal Distribution Function", "n . h", paste("D(n,h,",alpha,")"), d, alpha)

k <- 0.1
create_plot3D("Geometry Function", "n . h", paste("D(n,h,",k,")"), g, k)