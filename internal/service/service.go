package service

//TODO:
// When recursively templating:
//   - First check if we have DAG.
//   - If any cycles is spotted we should not allow such an operation to be performed.
// Additionally, we should ensure that no v1alpha1.Profile Custom Resource can be created if there are cycles:
//   - We need to create a Validating webhook that checks for cycles.
// Finally, we also need a form of runtime invalidation mechanism for dynamic non-DAGs:
//   - Indeed, to ensure that the content of "v1alpha1.ArbitraryResources"- or "v1alpha1.WebhookContent"-
//     v1alpha1.AdditionalContent does not contain cycles.
//   - We can create a DAG upon requesting all those information however, we should use BFS in order to avoid infinite
//     cycles. A max depth when running BFS might be a good solution.

// Template
//
// There is a body containing references to other configs that can themselves contain references to other configs.
//
// Templating happens in 3 phases:
//   1. find references
//   2. resolve references
//   3. render the template
//
// After rendering the template, we can recursively search for any references to a template.
