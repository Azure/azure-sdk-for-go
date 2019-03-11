// Package collection seeks to provide an expressive and readable way of working with basic data structures in Go.
//
// As a former .NET developer, I deeply missed writing programs in the style of Linq. Doing so enables concurrent/
// parallel reactive programs to be written in a snap. Go's functional nature enables us to have a very similar, if more
// verbose, experience.
//
// Take for example the scenario of printing the number of Go source files in a directory. Using this package,
// this takes only a few lines:
//
// 	myDir := collection.Directory{
// 		Location: "./",
// 	}
//
// 	results := myDir.Enumerate(nil).Where(func(x interface{}) bool {
// 		return strings.HasSuffix(x.(string), ".go")
// 	})
//
// 	fmt.Println(results.CountAll())
//
// A directory is a collection of filesystem entries, so we're able to iterate through them using the "Enumerate"
// function. From there, we filter on only file names that end with ".go". Finally, we print the number of entries that
// were encountered.
//
// This is a trivial example, but imagine building more elaborate pipelines. Maybe take advantage of the
// `SelectParallel` function which allows multiple goroutines to process a single transform at once, with their results
// being funnelled into the next phase of the pipeline. Suddenly, injecting new steps can be transparent.
package collection
