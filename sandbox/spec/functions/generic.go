package functions

// DotProduct generic function computes dot product for two vectors.
// Vectors should be of the eqal length
func DotProduct[F ~float32 | ~float64](v1, v2 []F) F {
	// the product `x * y` and the addition `s += x * y` are computed with
	// float32 or float64 precision, respectively, depending on the type argument for `F`
	var s F
	for i, x := range v1 {
		y := v2[i]
		s += x * y
	}
	return s
}
