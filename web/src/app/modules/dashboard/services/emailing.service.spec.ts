import { TestBed } from '@angular/core/testing';

import { EmailingService } from './emailing.service';

describe('EmailingService', () => {
  let service: EmailingService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EmailingService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
