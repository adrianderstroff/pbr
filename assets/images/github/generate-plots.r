library(plotly)
library(orca)

plotSurfaceContour <- function(fname, title, xlab, ylab, zlab, f, k) {
  normalize <- function(x)    { return(x/100.0)  }
  lambda    <- function(x, y) { return(f(x,y,k)) }
  
  xs <- sapply(c(0:100), normalize)
  ys <- sapply(c(0:100), normalize)
  zs <- outer(xs, ys, lambda)
  
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
  
  orca(fig, file = fname)
}

# geometry function
geomSmith <- function(a, k)      { return(a / (a * (1 - k) + k))               }
g         <- function(a1, a2, k) { return(geomSmith(a1, k) * geomSmith(a2, k)) }

# print the different functions
plotSurfaceContour("G00.png","Geometry Function (k=0.0)", "n * v", "n * l", "G(n,v,l)", g, 0.0)
plotSurfaceContour("G01.png","Geometry Function (k=0.1)", "n * v", "n * l", "G(n,v,l)", g, 0.1)
plotSurfaceContour("G05.png","Geometry Function (k=0.5)", "n * v", "n * l", "G(n,v,l)", g, 0.5)
plotSurfaceContour("G09.png","Geometry Function (k=0.9)", "n * v", "n * l", "G(n,v,l)", g, 0.9)
plotSurfaceContour("G10.png","Geometry Function (k=1.0)", "n * v", "n * l", "G(n,v,l)", g, 1.0)