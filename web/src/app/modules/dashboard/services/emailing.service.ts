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
    return this.httpClient.post(`${this.API_URL}/api/send-email`, { from, to, subject, body, tailoredResumeId });
  }

  getBulkEmailJobStatus$(jobId: string) {
    return this.httpClient.get<any>(`${this.API_URL}/api/send-email/jobs/${jobId}`);
  }

  generateAiDraftColdEmail$(from: string, to: string, companyName: string, jobDescription: string, templateType: string | null, jobUrls: string[], tailoredResumeId?: string): Observable<AiDraftColdMail> {
    const payload = { from, to, companyName, jobDescription, templateType, jobUrls, tailoredResumeId };
    return this.httpClient.post<AiDraftColdMail>(`${this.API_URL}/api/draft-with-ai`, payload);
  }

  getReferralEmailByUuid$(uuid: string): Observable<ReferralMailbox[]> {
    return this.httpClient.get<ReferralMailbox[]>(`${this.API_URL}/api/sent-referrals?uuid=${uuid}`);
  }

  pollReferralMailbox$(email: string, companyName?: string): Observable<ReferralMailbox[]> {
    let url = `${this.API_URL}/api/sent-referrals?email=${email}`;
    if (companyName && companyName.trim().length > 0) {
      url += `&company=${companyName}`;
    }
    return this.httpClient.get<ReferralMailbox[]>(url);
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
  tailoredResumeId: string | undefined;
  createdAt: Date;
}
