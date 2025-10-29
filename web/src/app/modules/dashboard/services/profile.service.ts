import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../../../environments/environment';
import { Observable } from 'rxjs';

export type ProfileInformation = {
  firstName: string;
  lastName: string;
  resume: File;
  email: string;
  about: string;
  country: string;
  notifications: {
    offers: boolean;
    pushNotifications: string;
    receiveEmails: boolean
  }
}

export interface ApiProfileInformation {
  id: string;
  firstName: string;
  lastName: string;
  resume: string;
  about: string;
  country: string;
  notifications: Notifications;
  email: string;
  profileSummary: string;
  extractedContent: string;
}

export interface Notifications {
  offers: boolean;
  pushNotifications: string;
  receiveEmails: boolean;
}

export interface ProfileAnalytics {
  totalEmails: number;
  companies: Company[];
  tailoredResumeCount: number;
  referralEmailCount: number;
}

export interface Company {
  companyName: string;
  totalEmails: number;
  distinctUsersCount: number;
}



@Injectable({
  providedIn: 'root'
})
export class ProfileService {
  private API_URL = environment.NG_REFERRER_BACKEND_API_URL;
  constructor(private readonly httpClient: HttpClient) { }

  get ownerEmailAddress(): string {
    // Returns the owner email address of the service
    return environment.ownerEmailAddress;
  }

  getProfileInformation$(email: string): Observable<ApiProfileInformation> {
    const headers = new HttpHeaders();
    headers.append("Content-Type", "multipart/form-data");

    return this.httpClient.get<ApiProfileInformation>(`${environment.NG_REFERRER_BACKEND_API_URL}/api/profile?email=${email}`, { headers })
  }

  searchPeople$(query: string): Observable<{ email: string, companyName: string }[]> {
    const headers = new HttpHeaders();
    headers.append("Content-Type", "application/json");

    return this.httpClient.get<{ email: string, companyName: string }[]>(`${this.API_URL}/api/profile/search-people?query=${query}`, { headers: headers, withCredentials: false })
  }

  updateProfileInformation$(profileInfo: ProfileInformation): Observable<ApiProfileInformation> {
    const headers = new HttpHeaders();
    headers.append("Content-Type", "multipart/form-data");

    // Create the form data for the multipart form data
    const formData = new FormData();
    formData.append("firstName", profileInfo.firstName);
    formData.append("lastName", profileInfo.lastName);
    formData.append("resume", profileInfo.resume, profileInfo.resume.name);
    formData.append("email", profileInfo.email);
    formData.append("about", profileInfo.about);
    formData.append("country", profileInfo.country);
    formData.append("notifications", JSON.stringify(profileInfo.notifications));

    return this.httpClient.post<ApiProfileInformation>(`${this.API_URL}/api/profile/information`, formData, {
      headers,
      withCredentials: false
    })
  }

  patchProfileInformation$(profileInfo: Partial<ApiProfileInformation>): Observable<ApiProfileInformation> {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.httpClient.patch<ApiProfileInformation>(`${this.API_URL}/api/profile`, profileInfo, { headers });
  }

  getProfileAnalytics$(userEmail: string): Observable<ProfileAnalytics> {
    const headers = new HttpHeaders();
    headers.append("Content-Type", "application/json");
    return this.httpClient.get<ProfileAnalytics>(`${this.API_URL}/api/profile/analytics?email=${userEmail}`, { headers });
  }

  tailorResumeWithJobDescription$(jobDescription: string, userEmail: string, companyName: string, jobRole: string): Observable<{ id: string }> {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.httpClient.post<{ id: string }>(
      `${this.API_URL}/api/profile/tailor-resume`,
      { jobDescription, userEmail, companyName, jobRole },
      { headers }
    );
  }

  getTailoredResumeById$(id: string): Observable<any> {
    return this.httpClient.get<any>(`${this.API_URL}/api/profile/tailored-resume/${id}`);
  }

  updateTailoredResumeMarkdown$(id: string, resumeMarkdown: string) {
    return this.httpClient.patch<{ success: boolean }>(
      `${this.API_URL}/api/profile/tailored-resume`,
      { id, resumeMarkdown }
    );
  }

  downloadResumeAsPDF$(parsedResumeContent: string): Observable<Blob> {
    return this.httpClient.post<Blob>(`${this.API_URL}/api/profile/export-pdf`, { resumeContent: parsedResumeContent }, { responseType: 'blob' as any })
  }

  getLatestTailoredResumes$(userEmail: string, companyName?: string) {
    let url = `${this.API_URL}/api/profile/tailored-resumes?email=${userEmail}`;
    if (companyName && companyName.trim().length > 0) {
      url += `&companyName=${encodeURIComponent(companyName.trim())}`;
    }
    return this.httpClient.get<any[]>(url);
  }
}
