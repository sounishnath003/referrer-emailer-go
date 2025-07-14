import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ProfileService } from '../../services/profile.service';
import { MarkdownModule } from 'ngx-markdown';
import { NgIf, DatePipe } from '@angular/common';

@Component({
    selector: 'app-tailored-resume-view',
    standalone: true,
    imports: [NgIf, MarkdownModule, DatePipe],
    template: `
    <div class="max-w-2xl mx-auto mt-8 p-8 bg-white rounded-xl shadow-lg">
      <div *ngIf="loading" class="text-blue-600 font-semibold flex items-center gap-2">
        <span class="animate-spin inline-block w-5 h-5 border-2 border-blue-400 border-t-transparent rounded-full"></span>
        Loading tailored resume...
      </div>
      <div *ngIf="error" class="text-red-600 font-semibold">{{ error }}</div>
      <ng-container *ngIf="resume">
        <h2 class="text-2xl font-bold mb-4">Tailored Resume</h2>
        <div class="text-xs text-gray-400 mb-2">Created: {{ resume.createdAt | date:'medium' }}</div>
        <markdown [data]="resume.resumeMarkdown" class="bg-gray-50 p-6 rounded-lg text-base leading-relaxed shadow-md"></markdown>
      </ng-container>
    </div>
  `
})
export class TailoredResumeViewComponent implements OnInit {
    resume: any = null;
    loading = true;
    error: string | null = null;

    constructor(private route: ActivatedRoute, private profileService: ProfileService) { }

    ngOnInit() {
        const id = this.route.snapshot.paramMap.get('id');
        if (id) {
            this.profileService.getTailoredResumeById$(id).subscribe({
                next: (res) => {
                    this.resume = res;
                    this.loading = false;
                },
                error: (err) => {
                    this.error = err.error?.error || 'Failed to load tailored resume.';
                    this.loading = false;
                }
            });
        } else {
            this.error = 'No resume ID provided.';
            this.loading = false;
        }
    }
} 