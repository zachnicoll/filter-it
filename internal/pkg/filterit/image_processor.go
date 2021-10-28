package filterit

func applyFilters(filters []int) {
	// TODO: Use image magik to apply each of the filters to the image
}

func processMessage(id string) {
	// TODO: Fetch document form DynamoDb with respective id

	// TODO: Set document's progress attribute to PROCESSING

	// TODO: Fetch image from S3 based on the document's image attribute

	// TODO: Apply filters to image

	// TODO: Re-upload filtered image to S3

	// TODO: Write new filenname to document's image attribute

	// TODO: Set document progress attribute to DONE
}

func WatchQueue() {
	// TODO: In a loop, check if the SQS queue has a new message

	// TODO: If message, spin off a subroutine and process the message - processMessage(id)

	// TODO: Log custom metric to CloudWatch
}
