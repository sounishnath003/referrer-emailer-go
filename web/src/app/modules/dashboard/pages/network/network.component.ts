import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Subject, debounceTime, distinctUntilChanged, switchMap } from 'rxjs';
import { SubheroComponent } from '../shared/subhero/subhero.component';
import { NetworkService, Contact } from '../../services/network.service';
import { ProfileService } from '../../services/profile.service';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-network',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, SubheroComponent],
  templateUrl: './network.component.html',
  styleUrl: './network.component.css'
})
export class NetworkComponent implements OnInit {
  searchQuery: string = '';
  contacts: Contact[] = [];
  isLoading: boolean = false;
  searchSubject = new Subject<string>();

  showAddModal: boolean = false;
  showConnectModal: boolean = false;
  selectedContact: Contact | null = null;

  contactForm: FormGroup;
  ownerEmail: string = '';

  constructor(
    private networkService: NetworkService,
    private profileService: ProfileService,
    private fb: FormBuilder,
    private router: Router
  ) {
    this.ownerEmail = this.profileService.ownerEmailAddress;

    this.contactForm = this.fb.group({
      name: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      company: ['', Validators.required],
      role: ['', Validators.required],
      linkedin: [''],
      notes: ['']
    });

    this.searchSubject.pipe(
      debounceTime(300),
      distinctUntilChanged(),
      switchMap(query => {
        this.isLoading = true;
        return this.networkService.getContacts(this.ownerEmail, query);
      })
    ).subscribe({
      next: (results) => {
        this.contacts = results;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
        this.contacts = [];
      }
    });
  }

  ngOnInit() {
    this.loadContacts();
  }

  loadContacts() {
    this.isLoading = true;
    this.networkService.getContacts(this.ownerEmail, this.searchQuery).subscribe({
      next: (data) => {
        this.contacts = data;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
      }
    });
  }

  onSearch(query: string) {
    this.searchSubject.next(query);
  }

  openAddModal() {
    this.showAddModal = true;
    this.contactForm.reset();
  }

  closeAddModal() {
    this.showAddModal = false;
  }

  openConnectModal(contact: Contact) {
    this.selectedContact = contact;
    this.showConnectModal = true;
  }

  closeConnectModal() {
    this.showConnectModal = false;
    this.selectedContact = null;
  }

  navigateToDraftWithAi() {
    if (!this.selectedContact) return;
    this.router.navigate(['/dashboard/draft-with-ai'], {
      queryParams: {
        template: 'Ask for Job Opportunities',
        to: this.selectedContact.email,
        companyName: this.selectedContact.company
      }
    });
    this.closeConnectModal();
  }

  navigateToDirectEmail() {
    if (!this.selectedContact) return;
    this.router.navigate(['/dashboard/email-drafter'], {
      queryParams: {
        to: this.selectedContact.email
        // Email Drafter doesn't explicitly support companyName in the form usually, 
        // but passing it anyway is fine.
      }
    });
    this.closeConnectModal();
  }

  onSubmitContact() {
    if (this.contactForm.invalid) return;

    const newContact: Contact = {
      ...this.contactForm.value,
      ownerId: this.ownerEmail
    };

    this.networkService.addContact(newContact).subscribe({
      next: () => {
        this.closeAddModal();
        this.loadContacts();
      },
      error: (err) => console.error(err)
    });
  }

  deleteContact(id: string) {
    if (confirm('Are you sure you want to delete this contact?')) {
      this.networkService.deleteContact(id, this.ownerEmail).subscribe(() => {
        this.loadContacts();
      });
    }
  }

  syncContacts() {
    this.isLoading = true;
    this.networkService.syncContacts(this.ownerEmail).subscribe({
      next: (res) => {
        // window.alert(res.message);
        this.loadContacts();
      },
      error: (err) => {
        console.error(err);
        this.isLoading = false;
      }
    });
  }
}
