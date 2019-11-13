resource "aws_dynamodb_table" "process_table" {
  name             = var.dyamodb_name
  read_capacity    = var.dyamodb_read_capacity
  write_capacity   = var.dyamodb_write_capacity
  hash_key         = var.dyamodb_hash_key
  range_key        = var.dyamodb_range_key
  stream_enabled   = var.dyamodb_stream_enabled
  stream_view_type = var.dyamodb_stream_view_type
  dynamic "attribute" {
    for_each = var.dynamodb_table_attributes
    content {
      # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
      # which keys might be set in maps assigned here, so it has
      # produced a comprehensive set here. Consider simplifying
      # this after confirming which keys can be set in practice.

      name = attribute.value.name
      type = attribute.value.type
    }
  }

  #  ttl                         = "${var.dynamodb_table_ttl}"
  tags = var.tags
}
