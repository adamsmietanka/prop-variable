package main

import (
	"math"
)

func InterpolateY(X, Z []float64, Y [][]float64, x, z float64) float64{
	// binary search for x upper bound
	i, k := 0, len(X)
	for i < k {
		h := int(uint(i+k) >> 1)
		if X[h] < x {
			i = h + 1
		} else {
			k = h
		}
	}
	// binary search for z upper bound
	j, k := 0, len(Z)
	for j < k {
		h := int(uint(j+k) >> 1)
		if Z[h] < z {
			j = h + 1
		} else {
			k = h
		}
	}
	if x < 0 || i == len(X) {
		return 0
	}
	if z < 0 || j == len(Z) {
		return 0
	}
	if i == 0 && j == 0 {
		return 0
	}
	if i == 0 {
		w := (z - Z[j-1]) / (Z[j] - Z[j-1])
		return (1-w)*Y[j-1][0] + w*Y[j][0]
	}
	if j == 0 {
		w := (x - X[i-1]) / (X[i] - X[i-1])
		return (1-w)*Y[0][i-1] + w*Y[0][i]
	}
	w := (X[i] - X[i - 1]) * (Z[j] - Z[j - 1])
	f11 := Y[j - 1][i - 1] * (X[i] - x) * (Z[j] - z)
	f12 := Y[j - 1][i]     * (X[i] - x) * (z - Z[j - 1])
	f21 := Y[j][i - 1]     * (x - X[i - 1]) * (Z[j] - z)
	f22 := Y[j][i]         * (x - X[i - 1]) * (z - Z[j - 1])
	return (f11 + f12 + f21 + f22) / w
}

func InterpolateZ(X, Y []float64, Z [][]float64, x, y float64) float64{
	// binary search for x upper bound
	i, k := 0, len(X)
	for i < k {
		h := int(uint(i+k) >> 1)
		if X[h] < x {
			i = h + 1
		} else {
			k = h
		}
	}
	// regular mesh along Y makes it easy
	j := int(math.Ceil((y - 10) * float64(len(Y) - 1) / 50))
	if x < 0 || i == len(X) {
		return 0
	}
	if y < 10 || 60 < y {
		return 0
	}
	if i == 0 && j == 0 {
		return 0
	}
	if i == 0 {
		w := (y - Y[j-1]) / (Y[j] - Y[j-1])
		return (1-w)*Z[j-1][0] + w*Z[j][0]
	}
	if j == 0 {
		w := (x - X[i-1]) / (X[i] - X[i-1])
		return (1-w)*Z[0][i-1] + w*Z[0][i]
	}
	w := (X[i] - X[i - 1]) * (Y[j] - Y[j - 1])
	f11 := Z[j - 1][i - 1] * (X[i] - x) * (Y[j] - y)
	f12 := Z[j - 1][i]     * (X[i] - x) * (y - Y[j - 1])
	f21 := Z[j][i - 1]     * (x - X[i - 1]) * (Y[j] - y)
	f22 := Z[j][i]         * (x - X[i - 1]) * (y - Y[j - 1])
	return (f11 + f12 + f21 + f22) / w
}

func BarycentricY(X, Y []float64, Z [][]float64, x, z float64) float64{
	// binary search for x upper bound
	i, k := 0, len(X)
	for i < k {
		h := int(uint(i+k) >> 1)
		if X[h] < x {
			i = h + 1
		} else {
			k = h
		}
	}
	if x < 0 || i == len(X) {
		return 0
	}
	if x == 0 {
		i += 1
	}
	// binary search for y upper bound
	j, k := 0, len(Y)
	for j < k {
		h := int(uint(j+k) >> 1)
		w := (x - X[i-1]) / (X[i] - X[i-1])
		zBound := (1-w) * Z[h][i-1] + w * Z[h][i]
		if zBound < z {
			j = h + 1
		} else {
			k = h
		}
	}
	if j == 0 || j == len(Z) {
		return 0
	}
	w := (x - X[i-1]) / (X[i] - X[i-1])
	triangleBound := (1-w)*Z[j-1][i-1] + w*Z[j][i]
	if z < triangleBound {
		// lower triangle /_|
		A := (X[i] - X[i-1]) * (Z[j][i] - Z[j-1][i])
		w1 := ((Z[j-1][i-1] - Z[j-1][i]) * (x - X[i]) + (X[i] - X[i-1]) * (z - Z[j-1][i])) / A
		return Y[j-1] + w1 * (Y[j] - Y[j-1])
	} else {
		// upper triangle |/
		A := (Z[j-1][i-1] - Z[j][i-1]) * (X[i] - X[i-1])
		w2 := ((Z[j][i-1] - Z[j][i]) * (x - X[i-1]) + (X[i] - X[i-1]) * (z - Z[j][i-1])) / A
		return Y[j] - w2 * (Y[j] - Y[j-1])
	}
}

func BarycentricZ(X, Y []float64, Z [][]float64, x, y float64) float64{
	// binary search for x upper bound
	i, k := 0, len(X)
	for i < k {
		h := int(uint(i+k) >> 1)
		if X[h] < x {
			i = h + 1
		} else {
			k = h
		}
	}
	if x < 0 || i == len(X) {
		return 0
	}
	if x == 0 {
		i += 1
	}
	// binary search for y upper bound
	j, k := 0, len(Y)
	for j < k {
		h := int(uint(j+k) >> 1)
		if Y[h] < y {
			j = h + 1
		} else {
			k = h
		}
	}
	if y < 10 || j == len(Y) {
		return 0
	}
	if y == 10 {
		j += 1
	}
	w := (x - X[i-1]) / (X[i] - X[i-1])
	triangleBound := (1-w)*Y[j-1] + w*Y[j]
	if y < triangleBound {
		// lower triangle /_|
		w1 := (y - Y[j-1]) / (Y[j] - Y[j-1])
		w2 := (X[i] - x) / (X[i] - X[i-1])
		w3 := 1 - w1 - w2
		return Z[j][i] * w1 + Z[j-1][i-1] * w2 + Z[j-1][i] * w3
	} else {
		// upper triangle |/
		w1 := (x - X[i-1]) / (X[i] - X[i-1])
		w2 := (y - Y[j]) / (Y[j-1] - Y[j])
		w3 := 1 - w1 - w2
		return Z[j][i] * w1 + Z[j-1][i-1] * w2 + Z[j][i-1] * w3
	}
}