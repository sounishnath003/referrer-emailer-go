import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class EmailingService {
  private API_URL = environment.NG_REFERRER_BACKEND_API_URL;
  constructor(private readonly httpClient: HttpClient) { }

  sendEmail$(from: string, to: string[], subject: string, body: string, tailoredResumeId?: string) {
    const payload: any = {
      from,
      to,
      subject,
      body
    };
    if (tailoredResumeId) {
      payload.tailoredResumeId = tailoredResumeId;
    }
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

  generateAiDraftColdEmail$(from: string, to: string, companyName: string, jobDescription: string, templateType: string, jobUrls: string[]): Observable<AiDraftColdMail> {
    const headers = new HttpHeaders();
    headers.append("Content-Type", "application/json");

    const payload = {
      from,
      to,
      companyName,
      jobDescription,
      templateType,
      jobUrls
    };

    return this.httpClient.post<AiDraftColdMail>(`${environment.NG_REFERRER_BACKEND_API_URL}/api/draft-with-ai`, payload, { headers: headers });
  }

  getReferralEmailByUuid$(uuid: string) {
    return this.httpClient.get<ReferralMailbox[]>(`${environment.NG_REFERRER_BACKEND_API_URL}/api/sent-referrals?uuid=${uuid}`)
  }
}

export interface AiDraftColdMail {
  mailSubject: string;
  mailBody: string;
  tailoredResumeId?: string;
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
