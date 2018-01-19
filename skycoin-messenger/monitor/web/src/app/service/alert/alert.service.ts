import { Injectable } from '@angular/core';
import swal from 'sweetalert2';

@Injectable()
export class AlertService {

  constructor() { }
  error(message: string) {
    return swal('Error', message, 'error');
  }
}
