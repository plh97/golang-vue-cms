resource "vultr_instance" "instance" {
  plan              = "vc2-1c-1gb"
  region            = "nrt"
  os_id             = 1743
  firewall_group_id = "acbe8a72-6ad8-413e-ad71-a49db226cd11"
  hostname          = "vultr.guest"
  label             = "calendar make appointment instance"
  ssh_key_ids       = [var.ssh_key_id, vultr_ssh_key.my_ssh_key.id]

  # provisioner "local-exec" {
  #   command = "sh modify_ip.sh ${self.public_ip}"
  # }
}
