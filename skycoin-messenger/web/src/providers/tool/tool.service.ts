import { Injectable } from '@angular/core';
import { MdDialog } from '@angular/material';
import { AlertDialogComponent } from '../../components';

@Injectable()
export class ToolService {

  constructor(private dialog: MdDialog) { }

  alert(title: string, message: string = '', type: string = 'info') {
    const ref = this.dialog.open(
      AlertDialogComponent,
      {
        position: { top: '10%' },
        panelClass: 'alert-dialog-panel',
        backdropClass: 'alert-backdrop',
        width: '23rem'
      });
    ref.componentInstance.title = title;
    ref.componentInstance.message = message;
    ref.componentInstance.type = type;
    return ref.afterClosed();
  }
  ShowInfoAlert(title: string = 'Information', message: string = '') {
    return this.alert(title, message, 'info');
  }
  ShowSuccessAlert(title: string = 'Information', message: string = '') {
    return this.alert(title, message, 'success');
  }
  ShowWarningAlert(title: string = 'Information', message: string = '') {
    return this.alert(title, message, 'warning');
  }
  ShowDangerAlert(title: string = 'Information', message: string = '') {
    return this.alert(title, message, 'danger');
  }
}
