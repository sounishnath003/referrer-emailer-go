import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../../../environments/environment';

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

@Injectable({
  providedIn: 'root'
})
export class ProfileService {
  private API_URL = environment.NG_REFERRER_BACKEND_API_URL;
  constructor(private readonly httpClient: HttpClient) { }

  getProfileInformation$(email: string) {
    const headers = new HttpHeaders();
    headers.append("Content-Type", "multipart/form-data");

    return this.httpClient.get(`${environment.NG_REFERRER_BACKEND_API_URL}/api/profile?email=${email}`, { headers })
  }

  updateProfileInformation$(profileInfo: ProfileInformation) {
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

    return this.httpClient.post(`${this.API_URL}/api/profile/information`, formData, {
      headers,
      withCredentials: false
    })
  }
}
