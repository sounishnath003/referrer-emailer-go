import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../../../environments/environment';

export type ProfileInformation = {
  firstName: string;
  lastName: string;
  resume: File;
  email: string;
  about: string;
}

@Injectable({
  providedIn: 'root'
})
export class ProfileService {
  private API_URL = environment.NG_REFERRER_BACKEND_API_URL;
  constructor(private readonly httpClient: HttpClient) { }

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

    return this.httpClient.post(`${this.API_URL}/api/profile/information`, formData, {
      headers,
      withCredentials: false
    })
  }
}
