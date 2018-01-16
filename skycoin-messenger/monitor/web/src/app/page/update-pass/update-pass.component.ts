import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ApiService } from '../../service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-updatepass-page',
  templateUrl: 'update-pass.component.html',
  styleUrls: ['./update-pass.component.scss'],
  encapsulation: ViewEncapsulation.None
})

export class UpdatePassComponent implements OnInit {
  updateForm = new FormGroup({
    oldpass: new FormControl('', [Validators.required, Validators.minLength(4), Validators.maxLength(20)]),
    newpass: new FormControl('', [Validators.required, Validators.minLength(4), Validators.maxLength(20)]),
  });
  status = 0;
  constructor(private api: ApiService, private router: Router) { }

  ngOnInit() {
    this.api.checkLogin().subscribe(result => {
    });
  }
  init() {
    this.status = 0;
  }
  update(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    const data = new FormData();
    data.append('oldPass', this.updateForm.get('oldpass').value);
    data.append('newPass', this.updateForm.get('newpass').value);
    this.api.updatePass(data).subscribe(result => {
      if (result) {
        this.router.navigate([{ outlets: { user: ['login'] } }]);
      }
    }, err => {
      this.status = 1;
    });
  }
}
