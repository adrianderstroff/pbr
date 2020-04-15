library(plotly)
library(orca)

# plotting function
createPlot <- function(fname, title, xlab, ylab, f, a) {
  # calculate points
  xs <- lapply(c(0:100), function(x) x/100)
  ys <- lapply(xs, function(y) f(y,a))
  data <- data.frame(xs, ys)
  
  # draw line plot
  fig <- plot_ly(data, x = ~xs, y = ~ys, type = 'scatter', mode = 'lines')
  fig <- fig %>% layout(
    title = title,
      xaxis = list(title = xlab, range=c(0,1)),
      yaxis = list(title = ylab, range=c(0,1))
    )
  
  orca(fig, file = fname)
}

# constants
width  <- 700
height <- 500

# normal distribution function
d <- function(nDoth, alpha) {
  (alpha^2) / (pi * (nDoth^2 * (alpha^2 - 1) + 1)^2)
}

# print plots for different alpha
createPlot(paste("D00.png"), "Normal Distribution Function (alpha=0.0)", "n * h", "D(n, h, 0.0)", d, 0.0)
createPlot(paste("D01.png"), "Normal Distribution Function (alpha=0.1)", "n * h", "D(n, h, 0.1)", d, 0.1)
createPlot(paste("D05.png"), "Normal Distribution Function (alpha=0.5)", "n * h", "D(n, h, 0.5)", d, 0.5)
createPlot(paste("D09.png"), "Normal Distribution Function (alpha=0.9)", "n * h", "D(n, h, 0.9)", d, 0.9)
createPlot(paste("D10.png"), "Normal Distribution Function (alpha=1.0)", "n * h", "D(n, h, 1.0)", d, 1.0)

f <- function(hDotv, fo) {
  return(fo + (1.0 - fo) * (1 - hDotv)^5)
}

createPlot(paste("F00.png"), "Fresnel Function (F0=0.0)", "h * v", "D(n, h, 0.0)", f, 0.0)
createPlot(paste("F01.png"), "Fresnel Function (F0=0.1)", "h * v", "D(n, h, 0.1)", f, 0.1)
createPlot(paste("F05.png"), "Fresnel Function (F0=0.5)", "h * v", "D(n, h, 0.5)", f, 0.5)
createPlot(paste("F09.png"), "Fresnel Function (F0=0.9)", "h * v", "D(n, h, 0.9)", f, 0.9)
createPlot(paste("F10.png"), "Fresnel Function (F0=1.0)", "h * v", "D(n, h, 1.0)", f, 1.0)
