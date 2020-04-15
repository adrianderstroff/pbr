library(plotly)
library(orca)

#___________________________________________LINE_CHART___________________________________________#

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
  
  # use orca to save to file
  orca(fig, file = fname)
}

#_______________________________________________NDF______________________________________________#

# normal distribution function
d <- function(nDoth, alpha) {
  (alpha^2) / (pi * (nDoth^2 * (alpha^2 - 1) + 1)^2)
}

# print plots for different alpha
createPlot(paste("D00.png"), "Normal Distribution Function (alpha=0.0)", "n * h", "D(n, h, 0.0)", d, 0.0)
createPlot(paste("D02.png"), "Normal Distribution Function (alpha=0.2)", "n * h", "D(n, h, 0.2)", d, 0.2)
createPlot(paste("D04.png"), "Normal Distribution Function (alpha=0.4)", "n * h", "D(n, h, 0.4)", d, 0.4)
createPlot(paste("D06.png"), "Normal Distribution Function (alpha=0.6)", "n * h", "D(n, h, 0.6)", d, 0.6)
createPlot(paste("D08.png"), "Normal Distribution Function (alpha=0.8)", "n * h", "D(n, h, 0.8)", d, 0.8)
createPlot(paste("D10.png"), "Normal Distribution Function (alpha=1.0)", "n * h", "D(n, h, 1.0)", d, 1.0)

#_____________________________________________FRESNEL____________________________________________#

# fresnel function
f <- function(hDotv, fo) {
  return(fo + (1.0 - fo) * (1 - hDotv)^5)
}

createPlot(paste("F00.png"), "Fresnel Function (F0=0.0)", "h * v", "D(n, h, 0.0)", f, 0.0)
createPlot(paste("F02.png"), "Fresnel Function (F0=0.2)", "h * v", "D(n, h, 0.2)", f, 0.2)
createPlot(paste("F04.png"), "Fresnel Function (F0=0.4)", "h * v", "D(n, h, 0.4)", f, 0.4)
createPlot(paste("F06.png"), "Fresnel Function (F0=0.6)", "h * v", "D(n, h, 0.6)", f, 0.6)
createPlot(paste("F08.png"), "Fresnel Function (F0=0.8)", "h * v", "D(n, h, 0.8)", f, 0.8)
createPlot(paste("F10.png"), "Fresnel Function (F0=1.0)", "h * v", "D(n, h, 1.0)", f, 1.0)

#__________________________________________SURFACE_PLOT__________________________________________#

plotSurfaceContour <- function(fname, title, xlab, ylab, zlab, f, k) {
  # setup helper functions
  normalize <- function(x)    { return(x/100.0)  }
  lambda    <- function(x, y) { return(f(x,y,k)) }
  
  # setup arrays for x,y,z
  xs <- sapply(c(0:100), normalize)
  ys <- sapply(c(0:100), normalize)
  zs <- outer(xs, ys, lambda)
  
  # define the surface plot with contours on the x-y plane
  fig <- plot_ly(
    type = 'surface',
    showscale = FALSE,
    contours = list(z = list(show=TRUE, usecolormap=TRUE, project=list(z=TRUE))),
    x = ~xs, y = ~ys, z = ~zs
  )
  fig <- fig %>% layout(
    title = title,
    showlegend=FALSE,
    scene = list(
      camera = list(
        eye=list(x=1.6,y=-1.6,z=0.5)
      ),
      xaxis = list(title = xlab, range=c(0,1)),
      yaxis = list(title = ylab, range=c(0,1)),
      zaxis = list(title = zlab, range=c(0,1.2))
    ))
  
  # use orca to save to file
  orca(fig, file = fname)
}

#________________________________________GEOMETRY_FUNCTION_______________________________________#

# geometry function
geomSmith <- function(a, k)      { return(a / (a * (1 - k) + k))               }
g         <- function(a1, a2, k) { return(geomSmith(a1, k) * geomSmith(a2, k)) }

# print the different functions
plotSurfaceContour("G00.png","Geometry Function (k=0.0)", "n * v", "n * l", "G(n,v,l)", g, 0.0)
plotSurfaceContour("G02.png","Geometry Function (k=0.2)", "n * v", "n * l", "G(n,v,l)", g, 0.2)
plotSurfaceContour("G04.png","Geometry Function (k=0.4)", "n * v", "n * l", "G(n,v,l)", g, 0.4)
plotSurfaceContour("G06.png","Geometry Function (k=0.6)", "n * v", "n * l", "G(n,v,l)", g, 0.6)
plotSurfaceContour("G08.png","Geometry Function (k=0.8)", "n * v", "n * l", "G(n,v,l)", g, 0.8)
plotSurfaceContour("G10.png","Geometry Function (k=1.0)", "n * v", "n * l", "G(n,v,l)", g, 1.0)