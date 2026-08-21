(setq file-name (nth 0 command-line-args-left))
(setq mode (or (nth 1 command-line-args-left) "open"))
(setq iterations
      (string-to-number
       (or (nth 2 command-line-args-left) "10000")))

(defun bench-elapsed (function)
  (let ((start (current-time)))
    (funcall function)
    (float-time (time-subtract (current-time) start))))

(let ((open-seconds
       (bench-elapsed
        (lambda ()
          (find-file-literally file-name)))))

  (princ
   (format "file_bytes=%d\n"
           (nth 7 (file-attributes file-name))))

  (princ
   (format "buffer_chars=%d\n"
           (buffer-size)))

  (princ
   (format "open_seconds=%.9f\n"
           open-seconds))

  (cond
   ((string= mode "open")
    nil)

   ((string= mode "lookup")
    (let ((start (current-time))
          (checksum 0)
          (buffer-end (max 1 (- (point-max) (point-min)))))

      (dotimes (i iterations)
        (goto-char
         (+ (point-min)
            (% (* i 7919) buffer-end)))

        (setq checksum
              (+ checksum (line-number-at-pos)))

        (when (= (% i 100) 0)
          (setq checksum
                (+ checksum (- (point) (point-min))))))

      (let ((elapsed
             (float-time
              (time-subtract
               (current-time)
               start))))

        (princ "operation=lookup\n")

        (princ
         (format "iterations=%d\n"
                 iterations))

        (princ
         (format "elapsed_seconds=%.9f\n"
                 elapsed))

        (princ
         (format "ns_per_op=%.2f\n"
                 (* 1000000000.0
                    (/ elapsed iterations))))

        (princ
         (format "checksum=%d\n"
                 checksum)))))

   ((string= mode "insert")
    (goto-char
     (/ (+ (point-min) (point-max)) 2))

    (let ((start (current-time)))

      (dotimes (_ iterations)
        (insert-char ?x))

      (let ((elapsed
             (float-time
              (time-subtract
               (current-time)
               start))))

        (princ "operation=insert\n")

        (princ
         (format "iterations=%d\n"
                 iterations))

        (princ
         (format "elapsed_seconds=%.9f\n"
                 elapsed))

        (princ
         (format "ns_per_op=%.2f\n"
                 (* 1000000000.0
                    (/ elapsed iterations)))))))

   (t
    (princ
     (format "unknown mode=%s\n"
             mode))
    (kill-emacs 1))))

(kill-emacs 0)