import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ProfileService } from '../../services/profile.service';
import { DatePipe, NgFor, NgIf } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subject } from 'rxjs';
import { debounceTime } from 'rxjs/operators';

@Component({
  selector: 'app-tailored-resume-list',
  standalone: true,
  imports: [NgFor, NgIf, DatePipe, FormsModule],
  template: `
    <div class="w-full mt-8 p-8 bg-white rounded-xl shadow-lg">
      <h2 class="text-xl font-bold mb-4">Your Latest Tailored Resumes</h2>
      <div class="mb-4 flex flex-col sm:flex-row gap-2 items-center">
        <input type="text" [(ngModel)]="companyFilter" (ngModelChange)="onCompanyFilterChange()" placeholder="Search by company name..."
          class="w-full sm:w-64 p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-400" />
      </div>
      <div *ngIf="loading" class="text-blue-600 font-semibold">Loading...</div>
      <div *ngIf="error" class="text-red-600 font-semibold">{{ error }}</div>
      <ul *ngIf="resumes && resumes.length> 0" class="w-full">
        <li *ngFor="let resume of resumes" class="mb-4 border-b pb-2">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div class="flex flex-col gap-y-1">
              <div class="font-semibold text-blue-600">{{ resume.companyName }} - {{ resume.jobRole }}</div>
              <div class="text-xs text-gray-500">{{ resume.createdAt | date:'medium' }}</div>
              <div class="text-sm text-gray-700 truncate max-w-xs">{{ resume.jobDescription }}</div>
            </div>
            <button class="ml-0 sm:ml-4 mt-2 sm:mt-0 bg-blue-600 text-white px-4 py-1 rounded hover:bg-blue-700 transition" (click)="goToResume(resume.id)">Open</button>
          </div>
        </li>
      </ul>
      <div *ngIf="(!resumes || resumes.length === 0) && !loading && !error" class="text-red-500 text-center">
        No tailored resumes found
      </div>
    </div>
  `
})
export class TailoredResumeListComponent implements OnInit {
  resumes: any[] = [];
  loading = true;
  error: string | null = null;
  companyFilter: string = '';
  private companyFilterSubject = new Subject<string>();

  constructor(private profileService: ProfileService, private router: Router) { }

  ngOnInit() {
    this.companyFilterSubject.pipe(debounceTime(300)).subscribe(() => {
      this.fetchResumes();
    });
    this.fetchResumes();
  }

  fetchResumes() {
    const userEmail = this.profileService.ownerEmailAddress;
    this.loading = true;
    this.profileService.getLatestTailoredResumes$(userEmail, this.companyFilter).subscribe({
      next: (res) => {
        this.resumes = res;
        this.loading = false;
      },
      error: (err) => {
        this.error = err.error?.error || 'Failed to load tailored resumes.';
        this.loading = false;
      }
    });
  }

  onCompanyFilterChange() {
    this.companyFilterSubject.next(this.companyFilter);
  }

  goToResume(id: string) {
    this.router.navigate([`/dashboard/craft-resume/${id}`]);
  }
} 