import {Component, inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA} from "@angular/material/dialog";

@Component({
  selector: 'app-confirm',
  imports: [],
  templateUrl: './confirm.component.html',
  styleUrl: './confirm.component.css'
})
export class ConfirmDialogComponent implements OnInit {
  data = inject(MAT_DIALOG_DATA);

  ngOnInit(): void {
    console.log(this.data);
  }


}

