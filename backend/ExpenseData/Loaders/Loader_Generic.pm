#!/usr/bin/perl

package Loader_Generic;
use Moose;
extends 'Loader';

# build string formats:
# filename
sub BUILD
{
    my ($self) = @_;
    my @buildParts = split (';' ,$self->build_string);
    $self->setFileName($buildParts[0]);
}

1;

