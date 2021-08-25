package swaggerdoc

// swagger:operation PATCH /user/account/{account} user update_specific_user_fullname
//
// update specific user fullname
// ```
// curl -X PATCH 'https://localhost/user/account/{{account}}' \
// |-H 'Authorization: {{Authorization}}' \
// |-H 'X-CSRF-Token: {{X-CSRF-Token}}' \
// |-H 'Content-Type: application/json' \
// |-d '{
//     "fullname":""
// }' \
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
// - name: body
//   in: body
//   schema:
//     type: object
//     required:
//        - fullname
//     properties:
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
