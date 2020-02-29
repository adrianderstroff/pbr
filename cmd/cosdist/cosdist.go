package main

import (
	"math"

	"github.com/adrianderstroff/pbr/pkg/cgm"
	"github.com/go-gl/mathgl/mgl32"
)

// Ray primitive consisting of a position o and a direction dir.
type Ray struct {
	o   mgl32.Vec3
	dir mgl32.Vec3
}

// Sphere consists of a position o and a radius r.
type Sphere struct {
	o mgl32.Vec3
	r float32
}

// AABB is an axis aligned bounding box defined by the min and max point.
type AABB struct {
	min mgl32.Vec3
	max mgl32.Vec3
}

// Hit record storing the parameter t, intersection p and surface normal n.
type Hit struct {
	t float32
	p mgl32.Vec3
	n mgl32.Vec3
}

func isBetween(a, min, max float32) bool {
	return min <= a && a <= max
}

func vecMin(a, b mgl32.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{
		cgm.Min32(a.X(), b.X()),
		cgm.Min32(a.Y(), b.Y()),
		cgm.Min32(a.Z(), b.Z()),
	}
}

func vecMax(a, b mgl32.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{
		cgm.Max32(a.X(), b.X()),
		cgm.Max32(a.Y(), b.Y()),
		cgm.Max32(a.Z(), b.Z()),
	}
}

func vecMul(a, b mgl32.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{
		a.X() * b.X(),
		a.Y() * b.Y(),
		a.Z() * b.Z(),
	}
}

func min3(a, b, c float32) float32 {
	return cgm.Min32(a, cgm.Min32(b, c))
}

func max3(a, b, c float32) float32 {
	return cgm.Max32(a, cgm.Max32(b, c))
}

func intersectSphere(ray *Ray, sphere *Sphere, tmin, tmax float32) (Hit, bool) {
	// solve squared term
	oc := ray.o.Sub(sphere.o)
	a := ray.dir.Dot(ray.dir)
	b := oc.Dot(ray.dir)
	c := oc.Dot(oc) - sphere.r*sphere.r
	discriminant := b*b - a*c

	if discriminant > 0 {
		// get both solutions
		t1 := (-b - cgm.Sqrt32(discriminant)) / a
		t2 := (-b + cgm.Sqrt32(discriminant)) / a

		// check which one is inside
		t1inside := isBetween(t1, tmin, tmax)
		t2inside := isBetween(t2, tmin, tmax)

		// early return
		if !t1inside && !t2inside {
			return Hit{}, false
		}

		// determine the right parameter t
		t := t1
		if t1inside && t2inside {
			t = cgm.Min32(t1, t2)
		} else if t2inside {
			t = t2
		}

		p := ray.o.Add(ray.dir.Mul(t))

		// set hit record
		hit := Hit{
			t: t,
			p: p,
			n: p.Sub(sphere.o).Normalize(),
		}

		return hit, true
	}

	return Hit{}, false
}

func intersectAABB(ray *Ray, aabb *AABB) mgl32.Vec3 {
	d := mgl32.Vec3{
		-1 / ray.dir.X(),
		-1 / ray.dir.Y(),
		-1 / ray.dir.Z(),
	}

	tMin := vecMul(aabb.min.Sub(ray.o), d)
	tMax := vecMul(aabb.max.Sub(ray.o), d)
	t1 := vecMin(tMin, tMax)
	t2 := vecMax(tMin, tMax)
	tNear := max3(t1.X(), t1.Y(), t1.Z())
	tFar := min3(t2.X(), t2.Y(), t2.Z())

	//t := cgm.Min32(tNear, tFar)
	t := tNear
	if tNear < 0 {
		t = tFar
	}

	return ray.o.Add(ray.dir.Mul(t))
}

func cosineDistribution(hit *Hit, r1, r2, a float32) mgl32.Vec3 {
	// calculate tangent and binormal
	n := hit.n.Normalize()
	t := mgl32.Vec3{0, 1, 0}
	if n.Dot(mgl32.Vec3{0, 1, 0}) == 1 {
		t = mgl32.Vec3{1, 0, 0}
	}
	b := t.Cross(n)
	t = n.Cross(b)

	// spherical coordinates
	sinTheta := cgm.Sqrt32(r1)
	cosTheta := cgm.Sqrt32(1 - sinTheta*sinTheta)
	psi := 2 * math.Pi * r2
	v1 := n.Mul(cosTheta)
	v2 := t.Mul(sinTheta * cgm.Cos32(psi) * a)
	v3 := b.Mul(sinTheta * cgm.Sin32(psi) * a)

	dir := v1.Add(v2).Add(v3)
	return dir.Normalize()
}
