package template

const COLLECTION_CODE = `
<?php
namespace {{NAMESPACE}};

use Krecek\Database\Annotation\Entity;
use Krecek\Database\StoredCollection;

/**
 * {{READABLE_NAME}} collection.
 *
 * @method {{ENTITY_NAME}} current()
 *
 * @Entity({{ENTITY_NAME}}::class)
 */
class {{COLLECTION_NAME}} extends StoredCollection
{

}
`
