package template

const REPOSITORY_CODE = `
<?php

namespace {{NAMESPACE}};


use Krecek\Database\Annotation\Collection;
use Krecek\Database\StorageRepository;

/**
 * Class {{REPOSITORY_NAME}}
 * @package {{NAMESPACE}}
 *
 * @method {{ENTITY_NAME}} create()
 * @method {{ENTITY_NAME}} get($key)
 * @method {{ENTITY_NAME}} findOne($condition)
 * @method {{COLLECTION_NAME}} listAll()
 * @method {{COLLECTION_NAME}} findMany($condition)
 *
 * @Collection({{COLLECTION_NAME}}::class)
 */

class {{REPOSITORY_NAME}} extends StorageRepository
{

}
`
