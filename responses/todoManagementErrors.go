package responses

const (
	TodoMissing       HttpMessage = "Please provide the details of the todo item."
	ErrorCreateTodo   HttpMessage = "We couldn't create your todo. Please try again."
	ErrorUpdateTodo   HttpMessage = "We couldn't update your todo. Please try again."
	SuccessAddTodo    HttpMessage = "Your todo has been successfully added!"
	ErrorNoTodosFound HttpMessage = "No todo items found. Start by adding a new todo!"
	InvalidTodoID     HttpMessage = "The provided todo ID is not valid."
	ErrorTodoNotFound HttpMessage = "We couldn't find the todo item you're looking for."
)
