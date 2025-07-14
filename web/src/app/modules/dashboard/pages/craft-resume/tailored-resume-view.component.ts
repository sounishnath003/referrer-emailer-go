import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ProfileService } from '../../services/profile.service';
import { MarkdownModule } from 'ngx-markdown';
import { NgIf, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import html2pdf from 'html2pdf.js';

@Component({
    selector: 'app-tailored-resume-view',
    standalone: true,
    imports: [NgIf, MarkdownModule, DatePipe, FormsModule],
    templateUrl: './tailored-resume-view.component.html'
})
export class TailoredResumeViewComponent implements OnInit {
    resume: any = null;
    loading = true;
    error: string | null = null;
    editMode = false;
    editableMarkdown = '';

    constructor(private route: ActivatedRoute, private profileService: ProfileService) { }

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
    }

    downloadPDF() {
        const element = document.getElementById('resume-pdf-content');
        if (element) {
            html2pdf().from(element).set({
                margin: 0.5,
                filename: 'tailored-resume.pdf',
                html2canvas: { scale: 2 },
                jsPDF: { unit: 'in', format: 'letter', orientation: 'portrait' }
            }).save();
        }
    }
} 