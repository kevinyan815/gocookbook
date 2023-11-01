package main

func VariableLenNoInitialSlice(size int) []int {
   slice := make([]int, 0)
   for i := 0; i < size; i++ {
      slice = append(slice, i)
   }
   return slice
}

func VariableLenShorterSlice(size int) []int {
   slice := make([]int, 100)
   for i := 0; i < size; i++ {
      slice = append(slice, i)
   }
   return slice
}

// Recommanded append in this way
// first predfine slice capability to parameter specified size
// then make it length to 0 at the beginning
// and append item into it in iteration.
// This would not cause the array underlying this slice to expend
// that means not memory copy and paste opertion
// that's why its more efficient then other ways
func VariableLenLargerSlice(size int) []int {
   slice := make([]int, 0, size)
   for i := 0; i < size; i++ {
      slice = append(slice, i)
   }
   return slice
}
