package swaggerdoc

// swagger:operation POST /signin auth SignIn
//
// signin
//
// will return jwt token
// ```
// curl -X POST 'https://localhost/signin' \
// |-H 'X-CSRF-Token: {{X-CSRF-Token}}' \
// |-H 'Content-Type: application/json' \
// |-d '{
//     "account": "",
//     "password": ""
// }' \
// |--insecure
// ```
//
//
// ---
// parameters:
// - name: X-CSRF-Token
//   in: header
//   required: true
//   type: string
// - name: body
//   in: body
//   schema:
//     type: object
//     required:
//        - account
//        - password
//     properties:
//        account:
//          type: string
//        password:
//          type: string
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
//       example: {"data":{"token":"..."},"error":""}
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
