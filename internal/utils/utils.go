package utils

// Map is a generic function that receives an array and a function.
// It applies the function to each element of the array and returns a new array.
// Example:
//
//	arr := []int{1, 2, 3}
//	result := Map(arr, func(n int) int {
//		return n * 2
//	})
//	fmt.Println(result) // [2, 4, 6]
func Map[T any, U any](arr []T, f func(T) U) []U {
	result := make([]U, 0)

	for _, v := range arr {
		result = append(result, f(v))
	}

	return result
}
