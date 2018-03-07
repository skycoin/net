import { Component, OnInit, ViewEncapsulation, Inject } from '@angular/core';
import { ApiService } from '../../service/api/api.service';
import { AlertService } from '../../service/alert/alert.service';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog } from '@angular/material';

const NOUPGRADE = 'No Upgrade Available';
const UPGRADE = 'Upgrade Available';

@Component({
  selector: 'app-update-card',
  templateUrl: './update-card.component.html',
  styleUrls: ['./update-card.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class UpdateCardComponent implements OnInit {
  progressValue = 0;
  progressTask = null;
  checkProgressTask = null;
  updateTimeout = 0;
  updateStatus = NOUPGRADE;
  hasUpdate = false;
  nodeUrl = '';
  checkUrl = 'http://messenger.skycoin.net:8100/api/version';
  dialogRef: MatDialogRef<UpdateCardComponent>;
  constructor(
    private api: ApiService,
    private alert: AlertService,
    @Inject(MAT_DIALOG_DATA) public data: { version?: string, tag?: string },
    private dialog: MatDialog
  ) { }

  ngOnInit() {
    this.api.checkUpdate(this.checkUrl, this.data.tag, this.data.version).subscribe((res: Update) => {
      this.hasUpdate = res.Update;
      if (this.hasUpdate) {
        this.updateStatus = UPGRADE;
      } else {
        this.updateStatus = NOUPGRADE;
      }
    });
  }
  startDownload(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.progressValue = 1;
    this.progressTask = setInterval(() => {
      this.progressValue += 1;
      if (this.progressValue >= 100) {
        clearInterval(this.progressTask);
      }
    }, 5000);
    this.api.runNodeupdate(this.nodeUrl).subscribe(r => {
      if (r) {
        this.checkProgressTask = setInterval(() => {
          this.api.getNodeupdateProcess(this.nodeUrl).subscribe(result => {
            if (result) {
              this.progressValue = 100;
              clearInterval(this.checkProgressTask);
              clearInterval(this.progressTask);
              this.dialog.closeAll();
              this.alert.success('Please restart the program to complete the update.');
            } else {
              if (this.updateTimeout >= 600) {
                this.progressValue = 0;
                this.dialog.closeAll();
                this.alert.error('Update timeout, restart the program and try again.');
                clearInterval(this.checkProgressTask);
                clearInterval(this.progressTask);
              }
              this.updateTimeout += 1;
            }
          });
        }, 500);
      }
    }, err => {
      this.alert.error('Update the timeout, check the network connection and restart the program and try again.');
      this.progressValue = 0;
      clearInterval(this.progressTask);
      this.dialog.closeAll();
    });
  }
  getUpgradeStatus() {
    return false;
  }
}

export interface Update {
  Force?: boolean;
  Update?: boolean;
  Latest?: string;
}
