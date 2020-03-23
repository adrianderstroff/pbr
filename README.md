# pbr

Simple implementation of a PBR Renderer using opengl 4.3. This project aims to get an understanding of simple Cook-Torrance BRDFs. 

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
framework on top of opengl (here using my fork that specifically supports glfw 3.2 with OpenGL 4.3), **hdr** is a loader for high dynamic range images
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

<p align="center">
![Microfacet Model](https://github.com/adrianderstroff/pbr/tree/master/assets/images/github/microfacet-model.png)
</p>

An important component of PBR is the **microfacet model**. Basically each surface can be described as microscopic perfect mirrors (**microfacets**). Depending on how the microfacets are aligned to each other, different kinds of surfaces can be modelled. A microfacet is described by its normal **h**. Are the normals of all microfacets parallel to each other, then the surface is completetly smooth and behaves like a perfect mirror. Are however the microfacet normals distributed in all directions then the surface is really rough as light is scattered in all directions.

Another important property of PBR techniques is that materials have to obey the **conservation of energy**, meaning that reflected light has to have the same or less energy than prior to the intersection with the material. 

Now one possible model to render PBR is the **reflectance equation**:
<p align="center">
<img src="https://latex.codecogs.com/png.latex?L_0(\mathbf{v})&space;=&space;\int_\Omega&space;f(\mathbf{l},&space;\mathbf{v})&space;L_i(\mathbf{l})&space;(\mathbf{n}&space;\cdot&space;\mathbf{l})&space;d&space;\omega_i" title="L_0(\mathbf{v}) = \int_\Omega f(\mathbf{l}, \mathbf{v}) L_i(\mathbf{l}) (\mathbf{n} \cdot \mathbf{l}) d \omega_i" />
</p>

Here *L<sub>0</sub>* is the **radiance** of a small patch of the surface in direction of the viewer **v**. 

The Ω represents the **hemisphere** from where light can hit the surface patch. The hemisphere can be thought of as a unit half sphere and light coming from any point of that half sphere towards the surface patch. 

The function *f(**l**, **v**)* is called the **bidirectional reflectance distribution function** (short BRDF). It returns a value between 0 and 1 and describes how much of the light's radiance *L<sub>i</sub>* coming from **l** is reflected in the viewer's direction **v**. As the BRDF is a distribution, the sum of all function values for all combinations of ***l*** and ***v*** have to add up to 1.

The dot product *(**n** ⋅ **l**)* describes how much of the light illuminates the surface patch. If the light rays hit the surface patch from above, where the direction of the light rays is parallel to the surface patch normal, the patch is most lit. The more the direction of the light rays is tilted the less illuminated the surface patch will be with no light hitting the surface patch when the light direction is perpendicular to the surface patch normal. 

Lastly the *dω<sub>i</sub>* describes a small patch on the hemisphere from where the light ray is coming.

**TODO describe radiance**

## TODOs

- Make Shader Error Message more readable