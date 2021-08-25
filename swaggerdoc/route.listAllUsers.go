package swaggerdoc

// swagger:operation GET /users user list_all_users
//
// list all users
// ```
// curl -X GET 'https://localhost/users?paging=1&sorting=desc' \
// |-H 'Authorization: {{Authorization}}' \
// |--insecure
// ```
//
// ---
// parameters:
// - name: Authorization
//   in: header
//   required: true
//   type: string
// - name: X-CSRF-Token
//   in: header
//   required: true
//   type: string
// - name: paging
//   in: query
//   required: false
//   default: 1
// - name: sorting
//   in: query
//   required: false
//   default: asc
//   enum: [asc, desc]
// responses:
//   '200':
//     description: OK
//     schema:
//       type: object
//       properties:
//         data:
//           type: array
//         error:
//           type: string
//       example: {"data":{"total":1,"users":[{"account":"","fullname":"","created_at":"","updated_at":""}]},"error":""}
//   '400':
//     description: something wrong
//     schema:
//       type: object
//       properties:
//         data:
//           nullable: true
//         error:
//           type: string
//       example: {"data":null,"error":"the error detail..."}
