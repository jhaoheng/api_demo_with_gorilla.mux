package swaggerdoc

// swagger:operation GET /user/me user get_user_detail_info
//
// get user detail info
// ```
// curl -X GET 'https://localhost/user/me' \
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
