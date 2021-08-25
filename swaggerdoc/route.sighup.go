package swaggerdoc

// swagger:operation POST /signup auth Signup
//
// signup
//
// create an user
// ```
// curl -X POST 'https://localhost/signup' \
// |-H 'X-CSRF-Token: {{X-CSRF-Token}}' \
// |-H 'Content-Type: application/json' \
// |-d '{
//     "account": "",
//     "password": "",
//     "fullname": ""
// }' \
// |--insecure
// ```
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
//        - fullname
//     properties:
//        account:
//          type: string
//        password:
//          type: string
//        fullname:
//          type: string
//   required: true
// responses:
//   '200':
//     description: OK
//     schema:
//       type: object
//       properties:
//         data:
//           type: string
//         error:
//           type: string
//       example: {"data":"success","error":""}
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
