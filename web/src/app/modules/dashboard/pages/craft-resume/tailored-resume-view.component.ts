import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ProfileService } from '../../services/profile.service';
import { MarkdownModule } from 'ngx-markdown';
import { NgIf, DatePipe, } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { marked } from 'marked';

@Component({
    selector: 'app-tailored-resume-view',
    standalone: true,
    imports: [NgIf, MarkdownModule, DatePipe, FormsModule],
    templateUrl: './tailored-resume-view.component.html'
})
export class TailoredResumeViewComponent implements OnInit {
    resume: any = null;
    loading = true;
    downloadLoading = false;
    error: string | null = null;
    editMode = false;
    editableMarkdown = '';

    constructor(private route: ActivatedRoute, private profileService: ProfileService, private router: Router) { }

    ngOnInit() {
        const id = this.route.snapshot.paramMap.get('id');
        if (id) {
            this.profileService.getTailoredResumeById$(id).subscribe({
                next: (res) => {
                    this.resume = res;
                    this.editableMarkdown = res.resumeMarkdown;
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

    toggleEdit() {
        this.editMode = !this.editMode;
        // When toggling from edit mode to preview mode, update the original resume content and save to DB
        if (!this.editMode && this.resume) {
            this.resume.resumeMarkdown = this.editableMarkdown;
            this.loading = true;
            this.profileService.updateTailoredResumeMarkdown$(this.resume.id, this.editableMarkdown).subscribe({
                next: () => {
                    this.loading = false;
                },
                error: (err) => {
                    this.error = err.error?.error || 'Failed to save changes.';
                    this.loading = false;
                }
            });
        }
    }

    downloadPDF() {
        // Get the actual rendered resume content (with Tailwind/CSS applied)
        this.downloadLoading = true;
        const element = marked(this.resume.resumeMarkdown, { async: false, });
        if (!element) {
            this.error = 'Please try again!';
            this.downloadLoading = false;
            return;
        }

        this.profileService.downloadResumeAsPDF$(element).subscribe({
            next: (blob) => {
                this.downloadLoading = false;
                const url = window.URL.createObjectURL(blob);
                const link = document.createElement('a');
                link.href = url;
                link.download = `resume_${new Date().getTime()}.pdf`;
                link.click();
            },
            error: (err) => {
                this.error = err.error?.error || 'Failed to download resume as pdf.';
            }
        })
    }

    sendEmail() {
        if (this.resume) {
            this.router.navigate(['/dashboard/draft-with-ai'], {
                queryParams: {
                    template: this.resume.jobRole,
                    tailoredResumeId: this.resume.id,
                    companyName: this.resume.companyName,
                }
            });
        }
    }
}
