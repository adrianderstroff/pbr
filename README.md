# pbr

Simple implementation of a PBR Renderer using opengl 4.3. This project aims to get an understanding of the Cook-Torrance BRDF. 

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
<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/microfacet-model.png?raw=true" title="Microfacet Model" />
</p>

An important component of PBR is the **microfacet model**. Basically each surface can be described as microscopic perfect mirrors (**microfacets**). Depending on how the microfacets are aligned to each other, different kinds of surfaces can be modelled. A microfacet is described by its normal **h**. Are the normals of all microfacets parallel to each other, then the surface is completetly smooth and behaves like a perfect mirror. Are however the microfacet normals distributed in all directions then the surface is really rough as light is scattered in all directions.

Another important property of PBR techniques is that materials have to obey the **conservation of energy**, meaning that reflected light has to have the same or less energy than prior to the intersection with the material. 

Now one possible model to render PBR is the **reflectance equation**:
<p align="center">
<img src="https://latex.codecogs.com/png.latex?L_0(\mathbf{v})&space;=&space;\int_\Omega&space;f(\mathbf{l},&space;\mathbf{v})&space;L_i(\mathbf{l})&space;(\mathbf{n}&space;\cdot&space;\mathbf{l})&space;d&space;\omega_i" title="L_0(\mathbf{v}) = \int_\Omega f(\mathbf{l}, \mathbf{v}) L_i(\mathbf{l}) (\mathbf{n} \cdot \mathbf{l}) d \omega_i" />
</p>

Here *L<sub>0</sub>* is the sum of reflected **radiance** of a small patch of the surface in direction of the viewer **v**. 

The Ω represents the **hemisphere** from where light can hit the surface patch. The hemisphere can be thought of as a unit half sphere and light coming from any point of that half sphere towards the surface patch. 

The function *f(**l**, **v**)* is called the **bidirectional reflectance distribution function** (short BRDF). It returns a value between 0 and 1 and describes how much of the light's radiance *L<sub>i</sub>* coming from **l** is reflected in the viewer's direction **v**. As the BRDF is a distribution, the sum of all function values for all combinations of ***l*** and ***v*** have to add up to 1.

The dot product *(**n** ⋅ **l**)* describes how much of the light illuminates the surface patch. If the light rays hit the surface patch from above, where the direction of the light rays is parallel to the surface patch normal, the patch is most lit. The more the direction of the light rays is tilted the less illuminated the surface patch will be with no light hitting the surface patch when the light direction is perpendicular to the surface patch normal. 

Lastly the *dω<sub>i</sub>* describes a small patch on the hemisphere from where the light ray is coming.

### Radiometric Quantities

The following information is taken from the [pbr-book](http://www.pbr-book.org/3ed-2018/Color_and_Radiometry/Radiometry.html), for more in-depth information check their website.

PBR bases it's theory on radiometric quantities to describe the "brightness" of a light when interacting with different materials. As brightness is not a property that can be physically described radiance and irradiance are used instead. 

Let's start with ***Light Energy*** first. A light source emits photons of different wavelength. Each wavelength *λ* has a specific energy *Q = hc / λ*, with *h* being Planck's constant and *c* the speed of light.

Next based on the Light Energy, we are now interested how much energy is emited per time. This can be done by measuring how much energy passes through a region. This quantity is called the ***Light Flux*** *Φ* and is simply the Light Energy differentiated by time *t*, or in formula *Φ = dQ / dt*.

If we now also take area that is illuminated by the light source we can calculate how much photons per time instance hit or path through this area *A*. This quantity is called ***Irradiance*** *E* and describes the density of the Light Flux. The irradiance is calculated as  *E = dΦ / dA* or *E = d<sup>2</sup> Q / (dt dA)*. The density of Light Flux leaving an area *A* is sometimes referred to as ***Radiant Exitance***.

To define ***Radiance*** we first have to define the ***Solid Angle***. In 2D we can take an object and project it onto a unit circle by tracing two lines from the center of the circle to both "ends" of the shape. We then measure the arc length of the shapes projection which is the angle (in radians). If the whole circle is covered we have an angle of *2π*.
The Solid Angle is an extension to a unit sphere. Now we project an object onto the surface of the unit sphere getting a specific area covering the sphere. The Solid angle is measured in steradians (sr). The full sphere has a Solid Angle of 4π while a hemisphere, which is a half sphere, has a Solid Angle of 2π. The Solid Angle can be represented by a vector with a direction and a magnitude that represents the area. 
We will assume that we will observe an infinitisimally small area so the Solid Angle will simply be a normalized vector *ω* indicating a direction relative to the center of the unit sphere.

With the concept we can describe ***Radiance*** as the Irradiance with respect to a solid angle *ω*. Thus the Radiance *L* can be calculated as *L = dE' / dω* or *L = d<sup>3</sup> Q / (dt dA' dω)*. It's important to note, that Radiance measures the Irradiance with respect to the area *A'* which is *A* projected on the plane *P*. The plane *P* is orthogonal to the Solid Angle *ω*, or in other words *ω* coincides with the normal of *P*.

### Cook-Torrance BRDF

Most of this part is taken from [https://learnopengl.com/PBR/Theory](https://learnopengl.com/PBR/Theory). This tutorial and the following ones on the website focus on an approach that simplifies to evaluates the integral of the Reflectance Equation in a precomputation step to get decent fps.

However the idea of this project is to have an implementation that helps understanding the Cook-Torrance BRDF and also is easy to translate into my raytracing project, where it won't be possible to evaluate the integral ahead of time. Thus all computations will be carried out as is while still trying to follow the theory of the *learnopengl* website. 

#### The BRDF

The Cook-Torrance BRDF consists of two terms, a diffuse Lambert term and a specular Cook-Torrance term *f(l,v) = k<sub>d</sub> f<sub>lambert</sub> + k<sub>s</sub> f<sub>cook-torrance</sub>*. The coefficients *k<sub>d</sub>* and *k<sub>s</sub>* have to add up to 1 to obey the law of conservation of energy. Here the coefficents are vectors, so they have to add up 1 componentwise.

#### Diffuse Term

The diffuse part is constant and describes ideal diffuse material that scatters light in all directions of the hemisphere evenly. Thus the diffuse part is *f<sub>lambert</sub> = c / π*. Here *c* is the color the surface or sometimes called the ***albedo***. 

The division of *π* comes from the fact, that we have to adhere to the conservation of energy. A detailed explanation can be found in Rory's blogpost [http://www.rorydriscoll.com/2009/01/25/energy-conservation-in-games/](http://www.rorydriscoll.com/2009/01/25/energy-conservation-in-games/). 
In short we want to make sure that the inequality 

<p align="center">
<img src="https://latex.codecogs.com/png.latex?\int_\Omega&space;\mathbf{c}_d&space;L_i&space;cos(\theta)&space;\delta&space;\omega&space;\le&space;L_i" title="\int_\Omega \mathbf{c}_d L_i cos(\theta) \delta \omega \le L_i" />
</p>

holds true. To make things easier, we only taking the diffuse part of the BRDF into account and also keeping the view vector *v* constant, we only integrate over the outgoing rays *l*. As *c<sub>d</sub>* and *L<sub>i</sub>* are constant in terms of the outgoing ray *l* they can be written in front of the integral. Then *L<sub>i</sub>* can be divided on both sides of the inequality. So now we just need to integrate the cosine term over the hemisphere. Since its hard to integrate over the hemisphere we can instead integrate over the halfsphere in polar coordinates using two integrals with *φ = [0,2π]* and *θ = [0,π/2]*. After solving the integral we end up with *π c<sub>d</sub> <= 1*. The surface color *c* is defined in the range (0,0,0) to (1,1,1) thus we need to divide by *π* to fullfill the inequality and thus to obey the conservation of energy. 

#### Specular Term

The specular term is a bit more complicated. It consists of three functions ***D***, ***F***, ***G*** and a normalization factor as shown below:

<p align="center">
<img src="https://latex.codecogs.com/png.latex?\dpi{100}&space;\small&space;\frac{\mathbf{D}&space;\mathbf{F}&space;\mathbf{G}}{4&space;(\mathbf{l}&space;\cdot&space;\mathbf{n})(\mathbf{v}&space;\cdot&space;\mathbf{n})}" title="\small \frac{\mathbf{D} \mathbf{F} \mathbf{G}}{4 (\mathbf{l} \cdot \mathbf{n})(\mathbf{v} \cdot \mathbf{n})}" />
</p>

TODO: understand the normalization factor. Check the link http://www.pbr-book.org/3ed-2018/Reflection_Models/Microfacet_Models.html for potential explanation.

In the Disney paper, they mention that the normalization comes from microfacet derivation (https://disney-animation.s3.amazonaws.com/library/s2012_pbs_disney_brdf_notes_v2.pdf).

In the following all three functions are being discussed.

##### Normal Distribution Function

http://www.reedbeta.com/blog/hows-the-ndf-really-defined/

| | | |
|-|-|-|
|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/D00.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/D02.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/D04.png?raw=true" title="Microfacet Model" />|
|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/D06.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/D08.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/D10.png?raw=true" title="Microfacet Model" />|

##### Geometry Function *G*

| | | |
|-|-|-|
|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/G00.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/G02.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/G04.png?raw=true" title="Microfacet Model" />|
|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/G06.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/G08.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/G10.png?raw=true" title="Microfacet Model" />|

##### Fresnel Reflection *F*

| | | |
|-|-|-|
|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/F00.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/F02.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/F04.png?raw=true" title="Microfacet Model" />|
|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/F06.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/F08.png?raw=true" title="Microfacet Model" />|<img src="https://github.com/adrianderstroff/pbr/blob/master/assets/images/github/F10.png?raw=true" title="Microfacet Model" />|

## Acknowledgement

### Obj Models

The models *bunny.obj* and *dragon.obj* are popular reconstructed models of scans taken from the [Stanford 3D Scanning repository](http://graphics.stanford.edu/data/3Dscanrep/).

### Gun Model and Textures

The gun.obj and textures in *assets/images/textures/material/gun* are the work of Andrew Maximum and are taken from [https://www.artstation.com/artwork/3k2](https://www.artstation.com/artwork/3k2).

## TODOs

- Make Shader Error Message more readable