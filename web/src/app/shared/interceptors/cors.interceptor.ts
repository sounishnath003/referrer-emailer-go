import { HttpInterceptorFn } from '@angular/common/http';

export const corsInterceptor: HttpInterceptorFn = (req, next) => {
  const modifiedRequest = req.clone({
    setHeaders: {
      "Access-Control-Allow-Origin": "http://localhost:4200",
      "Referrer": "http://localhost:4200/",
      "Content-Type": "application/json",
    }
  })
  return next(modifiedRequest);
};
