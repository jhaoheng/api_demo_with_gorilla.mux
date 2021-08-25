package swaggerdoc

// swagger:operation GET /csrf GetCSRFToken
//
// To operate the apis, get CSRF token at beginning
// ```
// curl -X GET 'https://localhost/csrf' --insecure
// ```
//
// ---
// responses:
//   '200':
//     description: OK
//     headers:
//       X-CSRF-Token:
//         schema:
//           type: string
//         description: CSRF Token
//     schema:
//       type: object
//       properties:
//         data:
//           type: string
//         error:
//           type: string
//       example: {"data":"success","error":""}
