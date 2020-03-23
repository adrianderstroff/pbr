# pbr

Simple implementation of a PBR Renderer using opengl 4.3. This project aims to 
get an understanding of simple Cook-Torrance Brdfs. 

## Requirements
This project requires a GPU with OpenGL 4.3+ support.

The following dependencies depend on cgo. To make them work under Windows a 
compatible version of **mingw** is necessary. Information can be found 
[here](https://github.com/go-gl/glfw/issues/91). In my case I used 
*x86_64-7.2.0-posix-seh-rt_v5-rev1*. After installing the right version of 
**mingw** you can continue by installing the dependencies that follow next.

This project depends on **glfw** for creating a window and providing a rendering 
context, **go-gl/gl** for providing bindings to OpenGL, **go-gl/mathgl** 
provides vector and matrix math for OpenGL, **nuklear** is a gui
framework on top of opengl, **hdr** is a loader for high dynamic range images
that can have light intensities greater than 1 for all color channels.
```
go get -u github.com/go-gl/glfw/v3.2/glfw       // window and event handling
go get -u github.com/go-gl/gl/v4.3-core/gl      // opengl api
go get -u github.com/go-gl/mathgl/mgl32         // opengl vector math
go get -u github.com/adrianderstroff/nuklear/nk // gui framework
go get -u github.com/mdouchement/hdr            // hdr image loader
```
After getting all dependencies the project should work without any errors.

## Theory

### PBR

PBR stands for **Physically Based Rendering** and describes rendering techniques that model how light interacts with a material or medium in a physically plausible manner.

An important component of PBR is the **microfacet model**. Basically each surface can be described as microscopic perfect mirrors (**microfacets**). Depending on how the microfacets are aligned to each other, different kinds of surfaces can be modelled. A microfacet is described by its normal **h**. Are the normals of all microfacets parallel to each other, then the surface is completetly smooth and behaves like a perfect mirror. Are however the microfacet normals distributed in all directions then the surface is really rough as light is scattered in all directions.

Another important property of PBR techniques is that materials have to obey the **conservation of energy**, meaning that reflected light has to have the same or less energy than prior to the intersection with the material. 

![Reflectance Equation](https://latex.codecogs.com/gif.latex?L_0(\mathbf{v})&space;=&space;\int_\Omega&space;f(\mathbf{l},&space;\mathbf{v})&space;L_i(\mathbf{l})&space;(\mathbf{n}&space;\cdot&space;\mathbf{l})&space;d&space;\omega_i)


## TODOs

- Make Shader Error Message more readable