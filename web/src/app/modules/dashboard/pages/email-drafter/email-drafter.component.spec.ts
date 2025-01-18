import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EmailDrafterComponent } from './email-drafter.component';

describe('EmailDrafterComponent', () => {
  let component: EmailDrafterComponent;
  let fixture: ComponentFixture<EmailDrafterComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EmailDrafterComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EmailDrafterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
