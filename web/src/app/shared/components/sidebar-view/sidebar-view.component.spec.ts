import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SidebarViewComponent } from './sidebar-view.component';

describe('SidebarViewComponent', () => {
  let component: SidebarViewComponent;
  let fixture: ComponentFixture<SidebarViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SidebarViewComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SidebarViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
