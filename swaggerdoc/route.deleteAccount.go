package swaggerdoc

// swagger:operation DELETE /user/account/{account} user delete_user
//
// delete user
// ```
// curl -X DELETE 'https://localhost/user/account/{{account}}' \
// |-H 'Authorization: {{Authorization}}' \
// |-H 'X-CSRF-Token: {{X-CSRF-Token}}' \
// |-H 'Cookie: {{Cookie}}' \
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
// - name: account
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
