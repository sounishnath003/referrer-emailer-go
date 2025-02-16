import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class EmailingService {
  private API_URL = environment.NG_REFERRER_BACKEND_API_URL;
  constructor(private readonly httpClient: HttpClient) { }

  sendEmail$(from: string, to: string[], subject: string, body: string) {
    const payload = {
      from,
      to,
      subject,
      body
    };
    const httpHeaders = new HttpHeaders();
    httpHeaders.append('Content-Type', 'application/json')

    return this.httpClient.post(`${this.API_URL}/api/send-email`, payload, {
      headers: httpHeaders,
      withCredentials: false,
    })
  }
}
