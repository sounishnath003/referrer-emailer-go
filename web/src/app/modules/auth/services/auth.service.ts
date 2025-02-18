import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(private readonly httpClient: HttpClient) { }

  login$(email: string, password: string) {
    const loginPayload = { email, password };
    // configure headers
    const headers = this.configureHttpHeaders();
    return this.httpClient.post(`${environment.NG_REFERRER_BACKEND_API_URL}/api/auth/login`, loginPayload, { headers });
  }

  signup$(email: string, password: string) {
    const signupPayload = { email, password };

    // configure headers
    const headers = this.configureHttpHeaders();

    return this.httpClient.post(`${environment.NG_REFERRER_BACKEND_API_URL}/api/auth/signup`, signupPayload, { headers });
  }

  private configureHttpHeaders() {
    const headers = new HttpHeaders();
    headers.append("Content-Type", "application/json");
    return headers;
  }
}
