import { NgFor } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { RouterLink } from '@angular/router';
import { ProfileService } from '../../../services/profile.service';

@Component({
  selector: 'app-menu',
  imports: [RouterLink, NgFor],
  templateUrl: './menu.component.html',
  providers: [ProfileService],
  styleUrl: './menu.component.css'
})
export class MenuComponent implements OnInit {
  isMenuOpen: boolean = false;
  ownerEmailAddress: string = '';

  constructor(private readonly profileService: ProfileService) { }

  ngOnInit(): void {
    this.ownerEmailAddress = this.profileService.ownerEmailAddress;
  }

  toggleMenu() {
    this.isMenuOpen = !this.isMenuOpen;
  }
}
