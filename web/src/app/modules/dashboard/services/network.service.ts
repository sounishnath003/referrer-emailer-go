import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../../environments/environment';

export interface Contact {
  id?: string;
  ownerId: string;
  name: string;
  email: string;
  mobile?: string;
  company: string;
  role: string;
  linkedin?: string;
  notes?: string;
  createdAt?: string;
  updatedAt?: string;
}

@Injectable({
  providedIn: 'root'
})
export class NetworkService {
  private API_URL = environment.NG_REFERRER_BACKEND_API_URL;

  constructor(private http: HttpClient) { }

  getContacts(ownerEmail: string, query: string = ''): Observable<Contact[]> {
    return this.http.get<Contact[]>(`${this.API_URL}/api/network/contacts?email=${ownerEmail}&query=${query}`);
  }

  addContact(contact: Contact): Observable<Contact> {
    return this.http.post<Contact>(`${this.API_URL}/api/network/contacts`, contact);
  }

  deleteContact(id: string, ownerEmail: string): Observable<any> {
    return this.http.delete(`${this.API_URL}/api/network/contacts/${id}?email=${ownerEmail}`);
  }

  syncContacts(ownerEmail: string): Observable<{ message: string, count: number }> {
    return this.http.post<{ message: string, count: number }>(`${this.API_URL}/api/network/contacts/sync?email=${ownerEmail}`, {});
  }
}
