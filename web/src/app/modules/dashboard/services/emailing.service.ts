import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../../../environments/environment';
import { Observable, switchMap, timer } from 'rxjs';

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

  pollReferralMailbox$(userEmailAddress: string): Observable<ReferralMailbox[]> {
    const headers = new HttpHeaders();
    headers.append("Content-Type", "application/json");

    // return timer(0, 3000).pipe(
    //   switchMap(() => this.httpClient.get<ReferralMailbox[]>(`${environment.NG_REFERRER_BACKEND_API_URL}/api/sent-referrals`, {
    //     headers: headers,
    //     params: { email: userEmailAddress }
    //   }))
    // )

    return this.httpClient.get<ReferralMailbox[]>(`${environment.NG_REFERRER_BACKEND_API_URL}/api/sent-referrals`, {
      headers: headers,
      params: { email: userEmailAddress }
    })
  }

  getReferralEmailByUuid$(uuid: string) {
    return this.httpClient.get<ReferralMailbox[]>(`${environment.NG_REFERRER_BACKEND_API_URL}/api/sent-referrals?uuid=${uuid}`)
  }
}


export interface ReferralMailbox {
  id: string;
  uuid: string;
  from: string;
  to: string[];
  subject: string;
  body: string;
  createdAt: Date;
}
