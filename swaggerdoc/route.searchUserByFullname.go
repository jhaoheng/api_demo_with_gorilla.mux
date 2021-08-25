package swaggerdoc

// swagger:operation GET /user/fullname/{fullname} user search_user_by_fullname
//
// search user by fullname
// ```
// curl -X GET 'https://localhost/user/fullname/{{fullname}}' \
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
// - name: fullname
//   in: path
//   type: string
//   required: true
// responses:
//   '200':
//     description: OK
//     schema:
//       type: object
//       properties:
//         data:
//           type: object
//         error:
//           type: string
//       example: {"data":{"account":"","fullname":"","created_at":"","updated_at":""},"error":""}
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
